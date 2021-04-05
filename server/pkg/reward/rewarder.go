package reward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	RedditAwardSilver   = "gid_1"                                      // 100
	RedditAwardGold     = "gid_2"                                      // 500
	RedditAwardPlatinum = "gid_3"                                      //1800
	RedditAwardLove     = "award_5eac457f-ebac-449b-93a7-eb17b557f03c" // cheap 20 coin award
)

type Rewarder interface {
	Reward(redditID string) error
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

func (r *RedditGold) Reward(redditID string) error {
	return r.giveGold(redditID)
}

func (r *RedditGold) giveGold(redditID string) error {
	token, err := r.getToken()
	if err != nil {
		return err
	}

	thingID := fmt.Sprintf("t2_%s", redditID)

	body := url.Values{}
	body.Add("api_type", "json")
	body.Add("thing_id", thingID)
	body.Add("fullname", thingID)
	body.Add("is_anonymous", "false")
	body.Add("message", "Oh chimpancy that, you were given an award.")
	body.Add("gild_type", RedditAwardSilver) //TODO: Should be gold.

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://oauth.reddit.com/api/v1/gold/gild/%s", thingID),
		bytes.NewBufferString(body.Encode()),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}
	r.logger.Debug("OK", zap.String("status", resp.Status), zap.String("request", body.Encode()), zap.String("thing_id", thingID), zap.String("response", string(respBody)))

	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with response: %s", string(respBody))
	}
	return nil
}

func (r *RedditGold) getToken() (string, error) {

	if r.cachedToken != nil && r.cachedToken.ExpiresAt.After(time.Now().Add(time.Minute*5)) {
		return r.cachedToken.Token, nil
	}

	body := url.Values{}
	body.Add("grant_type", "password")
	body.Add("username", r.cfg.GiverUsername)
	body.Add("password", r.cfg.GiverPassword)

	requestBody := bytes.NewBufferString(body.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://www.reddit.com/api/v1/access_token", requestBody)
	if err != nil {
		return "", errors.Wrap(err, "failed to create token request")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("scrimpton-bot (by /u/warmans)"))
	req.SetBasicAuth(r.cfg.AppID, r.cfg.AppSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute token request")
	}
	defer resp.Body.Close()

	token := &token{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&token); err != nil {
		return "", errors.Wrap(err, "failed to decode response as JSON")
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

func (r *NoopRewards) Reward(redditID string) error {
	r.logger.Info("Noop reward given", zap.String("redditID", redditID))
	return nil
}
