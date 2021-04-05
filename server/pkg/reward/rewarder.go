package reward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type Rewarder interface {
	Reward(username string) error
}

type RedditGoldCfg struct {
	GiverUsername string
	GiverPassword string
	AppID         string
	AppSecret     string
}

func (c *RedditGoldCfg) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.GiverUsername, prefix, "reddit-gold-username", "scrimpton-bot", "Username of gold giver")
	flag.StringVarEnv(fs, &c.GiverPassword, prefix, "reddit-gold-password", "", "Password of gold giver")
	flag.StringVarEnv(fs, &c.AppID, prefix, "reddit-gold-app-id", "My2Lx5sgOUvD_A", "App ID")
	flag.StringVarEnv(fs, &c.AppSecret, prefix, "reddit-gold-app-secret", "", "App secret")
}

func NewRedditGoldRewarder(logger *zap.Logger, cfg RedditGoldCfg) Rewarder {
	return &RedditGold{logger: logger, cfg: cfg}
}

type token struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	ExpiresAt time.Time
}

type RedditGold struct {
	logger      *zap.Logger
	cfg         RedditGoldCfg
	cachedToken *token
}

func (r *RedditGold) Reward(username string) error {
	return r.giveGold(username)
}

func (r *RedditGold) giveGold(username string) error {
	token, err := r.getToken()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://oauth.reddit.com/api/v1/gold/give/%s", username), bytes.NewBufferString("months=1"))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r.logger.Debug("OK", zap.String("status", resp.Status), zap.String("response", string(respBody)))

	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with response: %s", string(respBody))
	}
	return nil
}

func (r *RedditGold) getToken() (string, error) {

	if r.cachedToken != nil && r.cachedToken.ExpiresAt.After(time.Now().Add(time.Minute*5)) {
		return r.cachedToken.Token, nil
	}

	requestBody := bytes.NewBufferString(
		fmt.Sprintf("grant_type=password&username=%s&password=%s", r.cfg.GiverUsername, r.cfg.GiverPassword),
	)

	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.SetBasicAuth(r.cfg.AppID, r.cfg.AppSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	token := &token{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&token); err != nil {
		return "", err
	}
	token.ExpiresAt = time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
	r.cachedToken = token

	return token.Token, nil
}

func NewNoopRewarder(logger *zap.Logger) Rewarder {
	return &NoopRewards{logger: logger}
}

type NoopRewards struct {
	logger *zap.Logger
}

func (r *NoopRewards) Reward(username string) error {
	r.logger.Info("Noop reward given", zap.String("username", username))
	return nil
}
