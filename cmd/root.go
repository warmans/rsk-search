package cmd

import (
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/cmd/data"
	"github.com/warmans/rsk-search/cmd/index"
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "rsk-search",
		Short: "root command",
	}

	// search index commands
	root.AddCommand(index.RootCmd())
	root.AddCommand(data.RootCmd())
	//root.AddCommand(ImportCmd())

	return root
}
