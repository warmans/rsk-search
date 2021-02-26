package cmd

import (
	"encoding/json"
	"fmt"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func ImportCmd() *cobra.Command {

	var inputDir string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import raw scraped data and process it ready for storage/indexing",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			logger.Info("Starting import of dir...", zap.String("path", inputDir))

			dirEntries, err := ioutil.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {
				isNewfile, err := handleFile(path.Join(inputDir, dirEntry.Name()), path.Join(outputDir, dirEntry.Name()))
				if err != nil {
					logger.Fatal("Failed to handle path", zap.String("path", dirEntry.Name()), zap.Error(err))
				}
				if isNewfile {
					logger.Info("new file created")
				} else {
					logger.Info("old file was updated. note: this is current just an overwrite operation")
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./script/pilkipedia-scraper/raw", "Path to raw scraped files")
	cmd.Flags().StringVarP(&outputDir, "output-path", "o", "./var/raw", "Path to output files")

	return cmd
}

func handleFile(inPath, outPath string) (bool, error) {

	isNewFile := false

	_, err := os.Stat(outPath)
	if os.IsNotExist(err) {
		isNewFile = true
	}

	oldFile, err := os.OpenFile(outPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return false, fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		if err := oldFile.Close(); err != nil {
			panic(err)
		}
	}()

	newFile, err := os.OpenFile(inPath, os.O_RDONLY, 0666)
	if err != nil {
		return false, fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			panic(err)
		}
	}()

	//no comparison to do, just copy the new file into place.
	if isNewFile {
		if _, err := io.Copy(oldFile, newFile); err != nil {
			return false, err
		}
		return isNewFile, nil
	}

	// files must be compared
	oldData, err := decodeFile(oldFile)
	if err != nil {
		return false, fmt.Errorf("failed to decode old file: %w", err)
	}

	newData, err := decodeFile(newFile)
	if err != nil {
		return false, fmt.Errorf("failed to decode new file: %w", err)
	}

	if len(oldData.Transcript) != len(newData.Transcript) {
		return false, fmt.Errorf("wrong number of rows in new file. You'll need to manually fix the old file")
	}

	// todo: now go though the old data and update it row by row to maintain IDs

	// just overwrite the old file with the new data
	if err := oldFile.Truncate(0); err != nil {
		return false, err
	}
	if _, err := oldFile.Seek(0, 0); err != nil {
		return false, err
	}

	enc := json.NewEncoder(oldFile)
	enc.SetIndent("  ", "  ")
	return isNewFile, enc.Encode(newData) // just dump the new data for now
}

func decodeFile(f *os.File) (*models.Episode, error) {
	episode := &models.Episode{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(episode); err != nil {
		return nil, err
	}
	return episode, nil
}
