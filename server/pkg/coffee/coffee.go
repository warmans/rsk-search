package coffee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
	"io"
	"net/http"
)

type Config struct {
	AccessToken           string
	SupporterSyncInterval int64
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.AccessToken, prefix, "coffee-access-token", "", "Buy me a coffee access token")
	flag.Int64VarEnv(fs, &c.SupporterSyncInterval, prefix, "supporter-check-interval-seconds", 3600, "attempt to sync supporters to local DB every N seconds")
}

func NewClient(cfg *Config) *Client {
	return &Client{httpClient: http.DefaultClient, accessToken: cfg.AccessToken}
}

type Client struct {
	httpClient  *http.Client
	accessToken string
}

func (c *Client) Supporters() (*SupporterList, error) {
	req, err := http.NewRequest(http.MethodGet, "https://developers.buymeacoffee.com/api/v1/supporters", &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("User-Agent", "scrimpton-bot")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return nil, err
	}
	out := &SupporterList{
		Data:  make([]Supporter, 0),
		Total: 0,
	}
	if err := json.NewDecoder(buffer).Decode(&out); err != nil {
		return nil, fmt.Errorf("%v (status: %s, resp: %s)", err, resp.Status, buffer.String())
	}
	return out, nil
}
