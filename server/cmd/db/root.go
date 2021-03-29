package db

import (
	"github.com/spf13/cobra"
)

type dbconfig struct {
	path string
}

var cfg = dbconfig{}

func RootCmd() *cobra.Command {
	index := &cobra.Command{
		Use:   "db",
		Short: "commands related to to the search index",
	}

	index.PersistentFlags().StringVarP(&cfg.path, "db-path", "p", "./var/rsk.sqlite", "Path to sqlite DB")

	index.AddCommand(LoadCmd())
	index.AddCommand(LoadTscriptCmd())
	index.AddCommand(CreateRwTestdataCmd())

	return index
}
