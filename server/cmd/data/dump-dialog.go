package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
	"regexp"
)

func DumpDialog() *cobra.Command {

	var inputDir string
	var filePattern string

	cmd := &cobra.Command{
		Use:   "dump-dialog",
		Short: "Dump dialog to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			fileRegex, err := regexp.Compile(filePattern)
			if err != nil {
				return err
			}

			sentenceClosedRegex, err := regexp.Compile(".*[^a-zA-Z0-9]+")
			if err != nil {
				return err
			}

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			dirEntries, err := os.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {

				if dirEntry.IsDir() || !fileRegex.MatchString(dirEntry.Name()) {
					continue
				}

				episode := &models.Transcript{}
				if err := util.WithReadJSONFileDecoder(path.Join(inputDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				for idx, dialog := range episode.Transcript {
					if dialog.Type != models.DialogTypeChat || dialog.Actor == "" {
						continue
					}
					tryAddFullStop := false
					if idx+1 >= len(episode.Transcript) || dialog.Actor != episode.Transcript[idx+1].Actor {
						tryAddFullStop = true
					}
					suffix := " "
					if tryAddFullStop && !sentenceClosedRegex.MatchString(dialog.Content) {
						suffix = ". "
					}
					if _, err := os.Stdout.WriteString(dialog.Content + suffix); err != nil {
						return fmt.Errorf("failed to write to stdout: %w", err)
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to JSON files")
	cmd.Flags().StringVarP(&filePattern, "file-pattern", "p", ".*", "only include files that match this regex")

	return cmd
}
