package oauth

import (
	"github.com/karlseguin/ccache/v2"
	"github.com/lithammer/shortuuid/v3"
	"time"
)

func NewCSRFCache() *CSRFTokenCache {
	return &CSRFTokenCache{
		tokens: ccache.New(ccache.Configure().MaxSize(100).ItemsToPrune(10)),
	}
}

type CSRFTokenCache struct {
	tokens *ccache.Cache
}

func (c *CSRFTokenCache) NewCSRFToken(payload string) string {
	token := shortuuid.New()
	c.tokens.Set(token, payload, time.Minute*5)
	return token
}

func (c *CSRFTokenCache) VerifyCSRFToken(token string) (string, bool) {
	item := c.tokens.Get(token)
	if item == nil {
		return "", false
	}
	return item.Value().(string), !item.Expired()
}
