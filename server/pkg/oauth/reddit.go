package oauth

import (
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

//const RedditReturnURI = "http://scrimpton.com/oauth/reddit/return"
const RedditReturnURI = "http://localhost:4200/oauth/reddit/return"
const RedditApplicationID = "PytL99OIbkuUKw"

type Cfg struct {
	Secret            string
	KarmaLimit        int64
	MinAccountAgeDays int64
}

func (c *Cfg) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Secret, prefix, "oauth-secret", "", "Reddit oauth secret")
	flag.Int64VarEnv(fs, &c.KarmaLimit, prefix, "oath-karma-limit", 0, "Reddit oauth secret")
	flag.Int64VarEnv(fs, &c.MinAccountAgeDays, prefix, "oath-account-minage", 0, "")
}

type Identity struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	HasVerifiedEmail bool   `json:"has_verified_email"`
	Icon             string `json:"icon_img"`
	IsSuspended      bool   `json:"is_suspended"`

	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	TotalKarma   int64 `json:"total_karma"`

	Over18          bool `json:"over_18"`
}
