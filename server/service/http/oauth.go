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
	router.Path("/oauth/discord/return").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DiscordReturnHandler)))
}

func (c *OauthService) RedditReturnHandler(resp http.ResponseWriter, req *http.Request) {

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

	bearerToken, err := c.getBearerToken("reddit", "https://www.reddit.com/api/v1/access_token", code, c.oauthCfg.RedditAppID, c.oauthCfg.RedditSecret, state)
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

	authorIdentity := &models.Identity{
		ID:   ident.ID,
		Name: ident.Name,
		Icon: ident.Icon,
	}
	encodedIdentity, err := json.Marshal(authorIdentity)
	if err != nil {
		returnParams.Add("error", "Failed to decode response. Please report this on the scrimpton bug tracker.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	author := &models.Author{
		Name:          ident.Name,
		Identity:      string(encodedIdentity),
		OauthProvider: "reddit",
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

	token, err := c.auth.NewJWTForIdentity(author, authorIdentity)
	if err != nil {
		c.logger.Error("failed to create token", zap.Error(err))
		returnParams.Add("error", "failed to create token")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	returnParams.Add("token", token)
	http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
}

func (c *OauthService) getBearerToken(provider string, tokenEndpoint string, code string, appID string, appSecret string, state string) (string, error) {

	reqBody := url.Values{}
	if provider == "discord" {
		reqBody.Set("client_id", appID)
		reqBody.Set("client_secret", appSecret)
		reqBody.Set("scope", "identity")
	}
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("code", code)
	reqBody.Set("redirect_uri", c.oauthCfg.ProviderReturnURL(provider))

	fmt.Println("URL", reqBody.Encode())

	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, bytes.NewBufferString(reqBody.Encode()))
	if err != nil {
		c.logger.Error("failed to create access token request", zap.Error(err))
		return "", fmt.Errorf("unknown error")
	}
	req.Header.Set("User-Agent", "scrimpton-bot (by /u/warmans)")
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	if provider == "reddit" {
		req.SetBasicAuth(appID, appSecret)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("failed to request bearer token", zap.Error(err))
		return "", fmt.Errorf("failed to request bearer token")
	}
	defer resp.Body.Close()

	respBody := &bytes.Buffer{}
	if _, err := io.Copy(respBody, resp.Body); err != nil {
		return "", fmt.Errorf("failed to copy response body")
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error(
			"token request failed",
			zap.Error(err),
			zap.String("code", code),
			zap.String("provider", provider),
			zap.String("status", resp.Status),
			zap.String("body", respBody.String()),
		)
		return "", fmt.Errorf("ouath responded with unexpected status: %s", resp.Status)
	}

	buff := &bytes.Buffer{}
	if _, err := io.Copy(buff, respBody); err != nil {
		return "", fmt.Errorf("failed read oauth response")
	}

	fmt.Println("RESP", buff.String())

	response := struct {
		Error       string `json:"error"`
		AccessToken string `json:"access_token"`
	}{}
	if err := json.NewDecoder(buff).Decode(&response); err != nil {
		c.logger.Error("failed to decode bearer token", zap.Error(err), zap.String("response", buff.String()), zap.String("code", code))
		return "", fmt.Errorf("failed to request token")
	}
	if response.Error != "" {
		c.logger.Error("bearer token was an error", zap.String("error", response.Error))
		return "", fmt.Errorf("failed to authorize: was the token already used")
	}

	return response.AccessToken, nil
}

func (c *OauthService) getRedditIdentity(bearerToken string) (*oauth.RedditIdentity, string, error) {

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

	ident := &oauth.RedditIdentity{}
	if err := json.Unmarshal(responseBytes, ident); err != nil {
		c.logger.Error("failed to decode identity", zap.Error(err))
		return nil, "", fmt.Errorf("failed to parse identity")
	}

	// "new reddit" user icons for some reason now have a html encoded path. So just always try and decode it.
	ident.Icon = html.UnescapeString(ident.Icon)

	return ident, string(responseBytes), nil
}

func (c *OauthService) getDiscordIdentity(bearerToken string) (*oauth.DiscordIdentity, string, error) {

	if bearerToken == "" {
		c.logger.Error("blank bearer token, authorize must have failed")
		return nil, "", fmt.Errorf("authorization failed")
	}

	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/users/@me", bytes.NewBufferString(""))
	if err != nil {
		c.logger.Error("failed to create identity request", zap.Error(err))
		return nil, "", fmt.Errorf("failed to request identity")
	}

	req.Header.Set("User-Agent", "scrimpton-bot (by /u/warmans)")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", bearerToken))

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

	fmt.Println("TOKEN", buff.String())

	responseBytes := buff.Bytes()

	ident := &oauth.DiscordIdentity{}
	if err := json.Unmarshal(responseBytes, ident); err != nil {
		c.logger.Error("failed to decode identity", zap.Error(err))
		return nil, "", fmt.Errorf("failed to parse identity")
	}

	return ident, string(responseBytes), nil
}

func (c *OauthService) DiscordReturnHandler(resp http.ResponseWriter, req *http.Request) {
	errMessage := req.URL.Query().Get("error")
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	returnURL := fmt.Sprintf("%s%s/search", c.serviceConfig.Scheme, c.serviceConfig.Hostname)
	returnParams := url.Values{}

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

	bearerToken, err := c.getBearerToken("discord", "https://discord.com/api/oauth2/token", code, c.oauthCfg.DiscordAppID, c.oauthCfg.DiscordSecret, state)
	if err != nil {
		returnParams.Add("error", err.Error())
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	// reddit is surprisingly strict about how many requests you can send per second.
	time.Sleep(time.Second * 2)

	ident, _, err := c.getDiscordIdentity(bearerToken)
	if err != nil {
		returnParams.Add("error", err.Error())
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	authorIdentity := &models.Identity{
		ID:   ident.ID,
		Name: ident.Username,
		Icon: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.jpg", ident.ID, ident.Avatar),
	}
	encodedIdentity, err := json.Marshal(authorIdentity)
	if err != nil {
		returnParams.Add("error", "Failed to decode response. Please report this on the scrimpton bug tracker.")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	author := &models.Author{
		Name:          ident.Username,
		Identity:      string(encodedIdentity),
		OauthProvider: "discord",
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

	token, err := c.auth.NewJWTForIdentity(author, authorIdentity)
	if err != nil {
		c.logger.Error("failed to create token", zap.Error(err))
		returnParams.Add("error", "failed to create token")
		http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
		return
	}

	returnParams.Add("token", token)

	http.Redirect(resp, req, fmt.Sprintf("%s?%s", returnURL, returnParams.Encode()), http.StatusFound)
}
