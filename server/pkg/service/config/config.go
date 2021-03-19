package config

import (
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

type SearchServiceConfig struct {
	BleveIndexPath string
}

func (c *SearchServiceConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.BleveIndexPath, prefix, "bleve-index-path", "./var/rsk.bleve", "location of bleve search index")
}

