package bleve

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/search"
	"github.com/warmans/rsk-search/pkg/search/v1/index"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
)

func LoadCmd() *cobra.Command {

	var inputDir string

	cmd := &cobra.Command{
		Use:   "load",
		Short: "refresh the search index from the given directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			fmt.Printf("Using index %s...\n", indexCfg.path)

			rskIndex, err := bleve.Open(indexCfg.path)
			if err == bleve.ErrorIndexPathDoesNotExist {
				indexMapping, err := index.RskIndexMapping()
				if err != nil {
					logger.Fatal("failed to create mapping", zap.Error(err))
				}
				rskIndex, err = bleve.New(indexCfg.path, indexMapping)
				if err != nil {
					logger.Fatal("failed to create index", zap.Error(err))
				}
			}
			logger.Info("Populating index...")
			return populateIndex(inputDir, rskIndex, logger)
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw data files")

	return cmd
}

func populateIndex(inputDataPath string, idx bleve.Index, logger *zap.Logger) error {

	dirEntries, err := ioutil.ReadDir(inputDataPath)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		logger.Info("Parsing file...", zap.String("path", dirEntry.Name()))

		docs, err := documentsFromPath(path.Join(inputDataPath, dirEntry.Name()))
		if err != nil {
			return err
		}

		batch := idx.NewBatch()
		for _, d := range docs {
			if err := batch.Index(d.ID, d); err != nil {
				return err
			}
		}
		if err := idx.Batch(batch); err != nil {
			return err
		}
	}

	return nil
}

func documentsFromPath(filePath string) ([]search.DialogDocument, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer f.Close()

	episode := &models.Transcript{}

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(episode); err != nil {
		return nil, err
	}

	docs := []search.DialogDocument{}
	for _, v := range episode.Transcript {
		docs = append(docs, search.DialogDocument{
			ID:          v.ID,
			Mapping:     "dialog",
			Publication: episode.Publication,
			Series:      episode.Series,
			Episode:     episode.Episode,
			Date:        episode.ReleaseDate,
			ContentType: string(v.Type),
			Actor:       v.Actor,
			Position:    v.Position,
			Content:     v.Content,
		})
	}

	return docs, nil
}
