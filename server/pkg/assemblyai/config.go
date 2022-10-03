package assemblyai

import (
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

type Config struct {
	AccessToken string
}

func (c *Config) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.AccessToken, prefix, "assembly-ai-access-token", "", "Assembly AI API access token")
}
