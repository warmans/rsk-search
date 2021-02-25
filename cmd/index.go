package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/lithammer/shortuuid/v3"
	"github.com/spf13/cobra"
	"github.com/warmans/pilkipedia-scraper/pkg/models"
	"github.com/warmans/rsk-search/internal"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func IndexCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "index",
		Short: "refresh the search index from the given directory",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			fmt.Printf("Using index %s...\n", cfg.indexPath)

			rskIndex, err := bleve.Open(cfg.indexPath)
			if err == bleve.ErrorIndexPathDoesNotExist {
				indexMapping, err := internal.RskIndexMapping()
				if err != nil {
					logger.Fatal("failed to create mapping", zap.Error(err))
				}
				rskIndex, err = bleve.New(cfg.indexPath, indexMapping)
				if err != nil {
					logger.Fatal("failed to create index", zap.Error(err))
				}
			}

			logger.Info("Populating index...")
			return populateIndex(args[0], rskIndex, logger)

		},
	}
}

func populateIndex(inputDataPath string, idx bleve.Index, logger *zap.Logger) error {

	dirEntries, err := ioutil.ReadDir(inputDataPath)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {

		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		docs, err := documentsFromPath(path.Join(inputDataPath, dirEntry.Name()))
		if err != nil {
			return err
		}

		batch := idx.NewBatch()
		for _, d := range docs {
			if err := batch.Index(shortuuid.New(), d); err != nil {
				return err
			}
		}
		if err := idx.Batch(batch); err != nil {
			return err
		}
	}

	return nil
}

func documentsFromPath(filePath string) ([]internal.DialogDocument, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer f.Close()

	episode := &models.Episode{}

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(episode); err != nil {
		return nil, err
	}

	docs := []internal.DialogDocument{}
	for _, v := range episode.Transcript {
		docs = append(docs, internal.DialogDocument{
			Publication: episode.MetaValue(models.MetadataTypePublication),
			Series:      stringToIntOrZero(episode.MetaValue(models.MetadataTypeSeries)),
			Date:        episode.MetaValue(models.MetadataTypeDate),
			ContentType: string(v.Type),
			Actor:       v.Actor,
			Content:     v.Content,
		})
	}

	return docs, nil
}

func stringToIntOrZero(str string) int32 {
	i, _ := strconv.Atoi(str)
	return int32(i)
}
