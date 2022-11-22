package report

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "report",
		Short: "commands related to to the search index",
	}

	root.AddCommand(MonthlyRedditReport())

	return root
}
