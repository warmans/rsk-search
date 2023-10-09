package oauth

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

type Config struct {
	RedditAppID       string
	RedditSecret      string
	DiscordAppID      string
	DiscordSecret     string
	ReturnURL         string
	KarmaLimit        int64
	MinAccountAgeDays int64
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.RedditSecret, prefix, "oauth-reddit-secret", "", "reddit oauth secret")
	flag.StringVarEnv(fs, &c.RedditAppID, prefix, "oauth-reddit-app-id", "PytL99OIbkuUKw", "reddit application id")
	flag.StringVarEnv(fs, &c.DiscordAppID, prefix, "oauth-discord-app-id", "1161011766956937296", "discord application id")
	flag.StringVarEnv(fs, &c.DiscordSecret, prefix, "oauth-discord-secret", "", "discord application id")
	flag.StringVarEnv(fs, &c.ReturnURL, prefix, "oauth-return-url", "http://localhost:4200/oauth/%s/return", "return url must match reddit config")
	flag.Int64VarEnv(fs, &c.KarmaLimit, prefix, "oath-karma-limit", 10, "only allow accounts with at least this much karma")
	flag.Int64VarEnv(fs, &c.MinAccountAgeDays, prefix, "oath-account-minage", 1, "only allow accounts at least this many days old")
}

func (c *Config) ProviderReturnURL(provider string) string {
	return fmt.Sprintf(c.ReturnURL, provider)
}

type RedditIdentity struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	HasVerifiedEmail bool   `json:"has_verified_email"`
	Icon             string `json:"icon_img"`
	IsSuspended      bool   `json:"is_suspended"`

	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	TotalKarma int64 `json:"total_karma"`

	Over18 bool `json:"over_18"`
}

type DiscordIdentity struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
