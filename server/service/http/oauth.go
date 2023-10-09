package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/warmans/rsk-search/pkg/jwt"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/oauth"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"go.uber.org/zap"
	"html"
	"io"
	"net/http"
	"net/url"
	"time"
)

func NewOauthService(
	logger *zap.Logger,
	oauthCache *oauth.CSRFTokenCache,
	rwStore *rw.Conn,
	auth *jwt.Auth,
	oauthCfg *oauth.Config,
	serviceConfig config.SearchServiceConfig,
) *OauthService {
	return &OauthService{
		oauthCache:    oauthCache,
		logger:        logger.With(zap.String("component", "oauth-http-server")),
		rwStore:       rwStore,
		auth:          auth,
		oauthCfg:      oauthCfg,
		serviceConfig: serviceConfig,
	}
}

type OauthService struct {
	oauthCache    *oauth.CSRFTokenCache
	logger        *zap.Logger
	rwStore       *rw.Conn
	auth          *jwt.Auth
	oauthCfg      *oauth.Config
	serviceConfig config.SearchServiceConfig
}

func (c *OauthService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/oauth/reddit/return").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.RedditReturnHandler)))
}

func (c *OauthService) RedditReturnHandler(resp http.ResponseWriter, req *http.Request) {

	fmt.Println("URL", req.URL.String())

	returnURL := fmt.Sprintf("%s%s/search", c.serviceConfig.Scheme, c.serviceConfig.Hostname)
	returnParams := url.Values{}

	errMessage := req.URL.Query().Get("error")
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	if errMessage != "" {
		returnParams.Add("error", fmt.Sprintf("Auth failed with reason: %s", errMessage))
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	var ok bool
	returnURL, ok = c.oauthCache.VerifyCSRFToken(state)
	if !ok {
		returnParams.Add("error", "Request has invalid state, the request may have expired")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	bearerToken, err := c.getRedditBearerToken(code)
	if err != nil {
		returnParams.Add("error", err.Error())
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	// reddit is surprisingly strict about how many requests you can send per second.
	time.Sleep(time.Second * 2)

	ident, _, err := c.getRedditIdentity(bearerToken)
	if err != nil {
		returnParams.Add("error", err.Error())
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	// verify identity is in good standing
	if c.oauthCfg.KarmaLimit > 0 && ident.TotalKarma < c.oauthCfg.KarmaLimit {
		returnParams.Add("error", fmt.Sprintf("Account did not meet minimum karma requirements. You must have at least %d karma to authenticate.", c.oauthCfg.KarmaLimit))
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}
	if c.oauthCfg.MinAccountAgeDays > 0 {
		accountCreatedDate := time.Unix(int64(ident.CreatedUTC), 0)
		c.logger.Debug("Check account age", zap.Time("account_created", accountCreatedDate), zap.Time("current_time", time.Now().UTC()))
		if accountCreatedDate.After(time.Now().UTC().Add(0 - (time.Hour * 24 * time.Duration(c.oauthCfg.MinAccountAgeDays)))) {
			returnParams.Add("error", fmt.Sprintf("Account did not meet minimum age requirements. Your account must be at least %d days old to authenticate.", c.oauthCfg.MinAccountAgeDays))
			http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
			return
		}
	}
	if !ident.HasVerifiedEmail || ident.IsSuspended {
		returnParams.Add("error", "Account is unverified or suspended.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	encodedIdentity, err := json.Marshal(ident)
	if err != nil {
		returnParams.Add("error", "Failed to decode response. Please report this on the scrimpton bug tracker.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	author := &models.Author{
		Name:     ident.Name,
		Identity: string(encodedIdentity),
	}
	err = c.rwStore.WithStore(func(s *rw.Store) error {
		return s.UpsertAuthor(req.Context(), author)
	})
	if err != nil {
		c.logger.Error("failed to create local user", zap.Error(err))
		returnParams.Add("error", "failed to create local user")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	if author.Banned {
		returnParams.Add("error", "Account is not allowed.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
	}

	token, err := c.auth.NewJWTForIdentity(author, ident)
	if err != nil {
		c.logger.Error("failed to create token", zap.Error(err))
		returnParams.Add("error", "failed to create token")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	returnParams.Add("token", token)
	http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
}

func (c *OauthService) getRedditBearerToken(code string) (string, error) {
	requestBody := bytes.NewBufferString(
		fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, c.oauthCfg.ReturnURL),
	)

	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", requestBody)
	if err != nil {
		c.logger.Error("failed to create access token request", zap.Error(err))
		return "", fmt.Errorf("unknown error")
	}
	req.Header.Set("User-Agent", "scrimpton-bot (by /u/warmans)")
	req.SetBasicAuth(c.oauthCfg.RedditAppID, c.oauthCfg.RedditSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("failed to request bearer token", zap.Error(err))
		return "", fmt.Errorf("failed to request reddit token")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("reddit request failed", zap.Error(err), zap.String("code", code))
		return "", fmt.Errorf("reddit responded with unexpected status: %s", resp.Status)
	}

	buff := &bytes.Buffer{}
	if _, err := io.Copy(buff, resp.Body); err != nil {
		return "", fmt.Errorf("failed read reddit response")
	}

	response := struct {
		Error       string `json:"error"`
		AccessToken string `json:"access_token"`
	}{}
	if err := json.NewDecoder(buff).Decode(&response); err != nil {
		c.logger.Error("failed to decode bearer token", zap.Error(err), zap.String("response", buff.String()), zap.String("code", code))
		return "", fmt.Errorf("failed to request reddit token")
	}
	if response.Error != "" {
		c.logger.Error("bearer token was an error", zap.String("error", response.Error))
		return "", fmt.Errorf("failed to authorize:  was the token already used")
	}

	return response.AccessToken, nil
}

func (c *OauthService) getRedditIdentity(bearerToken string) (*oauth.Identity, string, error) {

	if bearerToken == "" {
		c.logger.Error("blank bearer token, authorize must have failed")
		return nil, "", fmt.Errorf("authorization failed")
	}

	req, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/v1/me", bytes.NewBufferString(""))
	if err != nil {
		c.logger.Error("failed to create identity request", zap.Error(err))
		return nil, "", fmt.Errorf("failed to request identity")
	}

	req.Header.Set("User-Agent", "scrimpton-bot (by /u/warmans)")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("failed to execute identity request", zap.Error(err))
		return nil, "", fmt.Errorf("failed to request identity")
	}
	defer resp.Body.Close()

	buff := &bytes.Buffer{}
	if _, err := buff.ReadFrom(resp.Body); err != nil {
		c.logger.Error("failed to copy response body", zap.Error(err))
		return nil, "", fmt.Errorf("failed to parse identity")
	}
	responseBytes := buff.Bytes()

	ident := &oauth.Identity{}
	if err := json.Unmarshal(responseBytes, ident); err != nil {
		c.logger.Error("failed to decode identity", zap.Error(err))
		return nil, "", fmt.Errorf("failed to parse identity")
	}

	// "new reddit" user icons for some reason now have a html encoded path. So just always try and decode it.
	ident.Icon = html.UnescapeString(ident.Icon)

	return ident, string(responseBytes), nil
}
