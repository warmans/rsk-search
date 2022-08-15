package coffee

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
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
	req, err := http.NewRequest(http.MethodGet, "https://developers.buymeacoffee.com/api/v1/supporters", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out := &SupporterList{
		Data:  make([]Supporter, 0),
		Total: 0,
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
