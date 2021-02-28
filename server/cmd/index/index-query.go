package index

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/internal"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/bleve_query"

	"go.uber.org/zap"
)

func QueryCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "query",
		Short: "query the index using the filter DSL",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			fmt.Printf("Using index %s...\n", indexCfg.path)

			rskIndex, err := bleve.Open(indexCfg.path)
			if err != nil {
				return err
			}

			filter, err := filter.Parse(args[0])
			if err != nil {
				return err
			}

			query, err := bleve_query.FilterToQuery(filter)
			if err != nil {
				return err
			}

			req := bleve.NewSearchRequest(query)
			req.Highlight = bleve.NewHighlightWithStyle("ansi")

			result, err := rskIndex.Search(req)
			if err != nil {
				return err
			}

			fmt.Println("HITS:")
			for _, v := range result.Hits {
				doc, err := rskIndex.Document(v.ID)
				if err != nil {
					return err
				}
				rawDoc := internal.DecodeDocument(doc)
				fmt.Printf("Score: %0.2f, Fragment: %s\n", v.Score, fmt.Sprint(v.Fragments))
				bytes, err := json.Marshal(rawDoc)
				if err != nil {
					return err
				}
				fmt.Println(string(bytes))
			}
			return nil
		},
	}

	return cmd
}
