package transcription

import (
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/tscript"
	"github.com/warmans/rsk-search/pkg/util"
	"os"
)

func ImportRawCmd() *cobra.Command {

	var episodeDir string
	var publication string
	var episode string

	cmd := &cobra.Command{
		Use:   "import-raw",
		Short: "use google cloud transcription service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return nil
			}

			var dialog []models.Dialog
			var synopsies []models.Synopsis
			err := util.WithExistingFile(args[0], func(f *os.File) error {
				var err error
				dialog, synopsies, err = tscript.Import(f)
				return err
			})
			if err != nil {
				return err
			}

			ep, err := data.LoadEpisode(episodeDir, publication, episode)
			if err != nil {
				return err
			}

			ep.Transcript = dialog
			ep.Synopsis = synopsies

			return data.ReplaceEpisodeFile(episodeDir, ep)
		},
	}

	cmd.Flags().StringVarP(&episodeDir, "episodes-dir", "d", "./var/data/episodes", "Path to raw data files")
	cmd.Flags().StringVarP(&publication, "publication", "p", "xfm", "")
	cmd.Flags().StringVarP(&episode, "episode", "e", "", "e.g. S2E01")

	return cmd
}
