package cmd

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/cmd/data"
	"github.com/warmans/rsk-search/cmd/db"
	"github.com/warmans/rsk-search/cmd/index"
	"github.com/warmans/rsk-search/cmd/server"
	"github.com/warmans/rsk-search/cmd/transcription"
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "rsk-search",
		Short: "root command",
	}

	// search index commands
	root.AddCommand(index.RootCmd())
	root.AddCommand(data.RootCmd())
	root.AddCommand(server.ServerCmd())
	root.AddCommand(db.RootCmd())
	root.AddCommand(transcription.RootCmd())

	return root
}
