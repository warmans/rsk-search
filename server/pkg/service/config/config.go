package config

import "github.com/warmans/rsk-search/pkg/flag"

type SearchServiceConfig struct {
	BleveIndexPath string
}

func (c *SearchServiceConfig) RegisterFlags(prefix string) {
	flag.StringVarEnv(&c.BleveIndexPath, prefix, "bleve-index-path", "./var/rsk.bleve", "location of bleve search index")
}

