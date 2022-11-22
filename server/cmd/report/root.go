package report

import (
	"github.com/spf13/cobra"
)

type dataConfig struct {
	dataDir  string
	audioDir string
}

var cfg = dataConfig{}

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "report",
		Short: "commands related to to the search index",
	}

	root.AddCommand(MonthlyRedditReport())

	return root
}
