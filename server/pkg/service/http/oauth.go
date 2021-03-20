package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/warmans/rsk-search/pkg/oauth"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewOauthService(logger *zap.Logger, oauthCache *oauth.CSRFTokenCache, redditOauthSecret string) *DownloadService {
	return &DownloadService{
		oauthCache:        oauthCache,
		redditOauthSecret: redditOauthSecret,
		logger:            logger.With(zap.String("component", "oauth")),
	}
}

type DownloadService struct {
	oauthCache        *oauth.CSRFTokenCache
	redditOauthSecret string
	logger            *zap.Logger
}

func (c *DownloadService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/oauth/reddit/return").Handler(http.HandlerFunc(c.RedditReturnHandler))
}

func (c *DownloadService) RedditReturnHandler(rw http.ResponseWriter, req *http.Request) {

	errMessage := req.URL.Query().Get("error")
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	if errMessage != "" {
		http.Error(rw, fmt.Sprintf("Auth failed with reason: %s", errMessage), http.StatusBadRequest)
		return
	}
	returnURL, ok := c.oauthCache.VerifyToken(state)
	if !ok {
		http.Error(rw, "Request has invalid state, the request may have expired", http.StatusBadRequest)
		return
	}

	bearerToken, err := c.getRedditBearerToken(code)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// reddit is surprisingly strict about how many requests you can send per second.
	time.Sleep(time.Second * 2)

	ident, err := c.getRedditIdentity(bearerToken)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Print(ident)

	//todo: store author in rw DB
	//todo: create JWT and attach to return URL

	http.Redirect(rw, req, returnURL, http.StatusFound)
}

func (c *DownloadService) getRedditBearerToken(code string) (string, error) {
	requestBody := bytes.NewBufferString(
		fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s", code, oauth.RedditReturnURI),
	)

	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", requestBody)
	if err != nil {
		c.logger.Error("failed to create access token request", zap.Error(err))
		return "", fmt.Errorf("unknown error")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.SetBasicAuth(oauth.RedditApplicationID, c.redditOauthSecret)

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

func (c *DownloadService) getRedditIdentity(bearerToken string) (*oauth.Identity, error) {

	if bearerToken == "" {
		c.logger.Error("blank bearer token, authorize must have failed")
		return nil, fmt.Errorf("authorization failed")
	}

	req, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/v1/me", bytes.NewBufferString(""))
	if err != nil {
		c.logger.Error("failed to create identity request", zap.Error(err))
		return nil, fmt.Errorf("failed to request identity")
	}

	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Error("failed to execute identity request", zap.Error(err))
		return nil, fmt.Errorf("failed to request identity")
	}
	defer resp.Body.Close()

	ident := &oauth.Identity{}
	if err := json.NewDecoder(resp.Body).Decode(&ident); err != nil {
		c.logger.Error("failed to decode identity", zap.Error(err))
		return nil, fmt.Errorf("failed to decode identity")
	}

	return ident, nil
}
