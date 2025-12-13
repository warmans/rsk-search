package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/transcript"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
)

func DumpPlaintext() *cobra.Command {

	var inputDir string
	var outputDir string
	var singleFile bool

	cmd := &cobra.Command{
		Use:   "dump-plaintext",
		Short: "Dump all transcripts in their plaintext format (instead of JSON).",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			var singleFileOutput io.WriteCloser
			if singleFile {
				var err error
				singleFileOutput, err = os.Create(path.Join(outputDir, "everything.txt"))
				if err != nil {
					return fmt.Errorf("failed to create single file output: %w", err)
				}
			}

			dirEntries, err := os.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {

				if dirEntry.IsDir() {
					continue
				}

				episode := &models.Transcript{}
				if err := util.WithReadJSONFileDecoder(path.Join(inputDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				rawTranscript, err := transcript.Export(episode.Transcript, episode.Synopsis, episode.Trivia, transcript.WithStripMetadata())
				if err != nil {
					return err
				}
				if singleFile {
					if _, err := fmt.Fprintf(singleFileOutput, "## %s\n\n", episode.ShortID()); err != nil {
						return err
					}
					if _, err := fmt.Fprint(singleFileOutput, rawTranscript, "\n\n"); err != nil {
						return err
					}
				} else {
					err = util.WithCreateOrReplaceFile(path.Join(outputDir, fmt.Sprintf("%s.txt", episode.ID())), func(f *os.File) error {
						_, err := f.WriteString(rawTranscript)
						if err != nil {
							return err
						}
						return nil
					})
				}
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to JSON files")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "./var/gen/plaintext", "Dump to this dir")
	cmd.Flags().BoolVarP(&singleFile, "single-file", "s", false, "Create a single file")

	return cmd
}
