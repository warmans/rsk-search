package config

import (
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

type SearchServiceConfig struct {
	BleveIndexPath  string
	Scheme          string
	Hostname        string
	RewardsDisabled bool
	FilesBasePath   string
}

func (c *SearchServiceConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Scheme, prefix, "scheme", "http://", "scheme to use for absolute links")
	flag.StringVarEnv(fs, &c.Hostname, prefix, "hostname", "localhost", "hostname to use for absolute links")
	flag.StringVarEnv(fs, &c.BleveIndexPath, prefix, "bleve-index-path", "./var/rsk.bleve", "location of bleve search index")
	flag.StringVarEnv(fs, &c.FilesBasePath, prefix, "files-base-path", "./var/data", "location files that can be downloaded")
	flag.BoolVarEnv(fs, &c.RewardsDisabled, prefix, "rewards-disabled", false, "Disable claiming rewards (but sill calculate them)")
}
