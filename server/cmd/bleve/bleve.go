package bleve

import (
	"github.com/spf13/cobra"
)

type indexConfig struct {
	path string
}

var indexCfg = indexConfig{}

func RootCmd() *cobra.Command {
	index := &cobra.Command{
		Use:   "bleve",
		Short: "commands related to to the search index",
	}

	index.PersistentFlags().StringVarP(&indexCfg.path, "index-path", "p", "./var/rsk.bleve", "Path to bleve index")

	index.AddCommand(LoadCmd())
	index.AddCommand(QueryCmd())

	return index
}
