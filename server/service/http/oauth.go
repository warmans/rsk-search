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
) *DownloadService {
	return &DownloadService{
		oauthCache:    oauthCache,
		logger:        logger.With(zap.String("component", "oauth-http-server")),
		rwStore:       rwStore,
		auth:          auth,
		oauthCfg:      oauthCfg,
		serviceConfig: serviceConfig,
	}
}

type DownloadService struct {
	oauthCache    *oauth.CSRFTokenCache
	logger        *zap.Logger
	rwStore       *rw.Conn
	auth          *jwt.Auth
	oauthCfg      *oauth.Config
	serviceConfig config.SearchServiceConfig
}

func (c *DownloadService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/oauth/reddit/return").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.RedditReturnHandler)))
}

func (c *DownloadService) RedditReturnHandler(resp http.ResponseWriter, req *http.Request) {

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

	ident, rawIdentityJSON, err := c.getRedditIdentity(bearerToken)
	if err != nil {
		returnParams.Add("error", err.Error())
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	// verify identity is in good standing
	if c.oauthCfg.KarmaLimit > 0 && ident.TotalKarma < c.oauthCfg.KarmaLimit {
		returnParams.Add("error", "Account did not meet minimum karma requirements.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}
	if c.oauthCfg.MinAccountAgeDays > 0 && time.Unix(int64(ident.CreatedUTC), 0).Add(0-(time.Hour*24*time.Duration(c.oauthCfg.MinAccountAgeDays))).Before(time.Now().UTC()) {
		returnParams.Add("error", "Account did not meet minimum age requirements.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}
	if !ident.HasVerifiedEmail || ident.IsSuspended {
		returnParams.Add("error", "Account is unverified or suspended.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	author := &models.Author{
		Name:     ident.Name,
		Identity: rawIdentityJSON,
	}
	err = c.rwStore.WithStore(func(s *rw.Store) error {
		return s.UpsertAuthor(req.Context(), author)
	})
	if err != nil {
		c.logger.Error("failed to create local user", zap.Error(err))
		returnParams.Add("error", fmt.Sprintf("failed to create local user"))
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	if author.Banned == true {
		returnParams.Add("error", "Account is not allowed.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
	}

	token, err := c.auth.NewJWTForIdentity(author, ident)
	if err != nil {
		c.logger.Error("failed to create token", zap.Error(err))
		returnParams.Add("error", fmt.Sprintf("failed to create token"))
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	returnParams.Add("token", token)
	http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
}

func (c *DownloadService) getRedditBearerToken(code string) (string, error) {
	requestBody := bytes.NewBufferString(
		fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, c.oauthCfg.ReturnURL),
	)

	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", requestBody)
	if err != nil {
		c.logger.Error("failed to create access token request", zap.Error(err))
		return "", fmt.Errorf("unknown error")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.SetBasicAuth(c.oauthCfg.AppID, c.oauthCfg.Secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("failed to request bearer token", zap.Error(err))
		return "", fmt.Errorf("failed to request reddit token")
	}
	defer resp.Body.Close()

	response := struct {
		Error       string `json:"error"`
		AccessToken string `json:"access_token"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("failed to decode bearer token", zap.Error(err))
		return "", fmt.Errorf("failed to request reddit token")
	}
	if response.Error != "" {
		c.logger.Error("bearer token was an error", zap.String("error", response.Error))
		return "", fmt.Errorf("failed to authorize:  was the token already used")
	}

	return response.AccessToken, nil
}

func (c *DownloadService) getRedditIdentity(bearerToken string) (*oauth.Identity, string, error) {

	if bearerToken == "" {
		c.logger.Error("blank bearer token, authorize must have failed")
		return nil, "", fmt.Errorf("authorization failed")
	}

	req, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/v1/me", bytes.NewBufferString(""))
	if err != nil {
		c.logger.Error("failed to create identity request", zap.Error(err))
		return nil, "", fmt.Errorf("failed to request identity")
	}

	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
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

	return ident, string(responseBytes), nil
}
