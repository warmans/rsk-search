package models

import (
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/util"
	"time"
)

type Changelog struct {
	Date    time.Time
	Content string
}

func (c *Changelog) Proto() *api.Changelog {
	return &api.Changelog{
		Date:    c.Date.Format(util.ShortDateFormat),
		Content: c.Content,
	}
}
