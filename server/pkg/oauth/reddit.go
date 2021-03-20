package oauth

import (
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

//const RedditReturnURI = "http://scrimpton.com/oauth/reddit/return"
const RedditReturnURI = "http://localhost:4200/oauth/reddit/return"
const RedditApplicationID = "PytL99OIbkuUKw"

type Cfg struct {
	Secret string
}

func (c *Cfg) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Secret, prefix, "oauth-secret", "", "Reddit oauth secret")
}

type Identity struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	HasVerifiedEmail bool   `json:"has_verified_email"`
	Icon             string `json:"icon_img"`
	IsSuspended      bool   `json:"is_suspended"`

	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`

	CommentKarma int64 `json:"comment_karma"`
	LinkKarma    int64 `json:"link_karma"`
	TotalKarma   int64 `json:"total_karma"`

	Over18          bool `json:"over_18"`
	PreferNightmode bool `json:"pref_nightmode"`
}
