module github.com/warmans/rsk-search

go 1.13

require (
	github.com/blevesearch/bleve/v2 v2.0.2
	github.com/blevesearch/bleve_index_api v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/lithammer/shortuuid/v3 v3.0.6
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.4.0
	github.com/warmans/pilkipedia-scraper v0.0.0
	go.uber.org/zap v1.10.0
)

replace github.com/warmans/pilkipedia-scraper => ../pilkipedia-scraper
