package cmd

import (
	"github.com/spf13/cobra"
)

type config struct {
	indexPath string
}

var cfg = config{}

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "rsk-search",
		Short: "root command",
	}

	root.PersistentFlags().StringVarP(&cfg.indexPath, "index-path", "i", "./var/rsk.bleve", "Path to bleve index")

	root.AddCommand(IndexCmd())
	root.AddCommand(QueryCmd())

	return root
}
