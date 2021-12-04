package blugeindex

import (
	"encoding/json"
	"fmt"
	"github.com/blugelabs/bluge"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/search"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
)

type indexConfig struct {
	path string
}

var indexCfg = indexConfig{}

func RootCmd() *cobra.Command {
	index := &cobra.Command{
		Use:   "bluge",
		Short: "commands related to to the bluge search index",
	}

	index.PersistentFlags().StringVarP(&indexCfg.path, "index-path", "p", "./var/rsk.bluge", "Path to bluge index")

	index.AddCommand(LoadCmd())

	return index
}

func LoadCmd() *cobra.Command {

	var inputDir string

	cmd := &cobra.Command{
		Use:   "load",
		Short: "refresh the search index from the given directory",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			fmt.Printf("Using index %s...\n", indexCfg.path)

			config := bluge.DefaultConfig(indexCfg.path)

			rskIndex, err := bluge.OpenWriter(config)
			if err != nil {
				return err
			}
			logger.Info("Populating index...")
			return populateIndex(inputDir, rskIndex, logger)
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw data files")

	return cmd
}

func populateIndex(inputDir string, writer *bluge.Writer, logger *zap.Logger) error {

	dirEntries, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		docs, err := documentsFromPath(path.Join(inputDir, dirEntry.Name()))
		if err != nil {
			return err
		}

		batch := bluge.NewBatch()
		for _, d := range docs {

			doc := bluge.NewDocument(d.ID)
			doc.AddField(bluge.NewKeywordField("publication", d.Publication).StoreValue())
			doc.AddField(bluge.NewNumericField("series", float64(d.Series)).StoreValue().Sortable())
			doc.AddField(bluge.NewDateTimeField("date", d.Date).Sortable())
			doc.AddField(bluge.NewKeywordField("actor", d.Actor).StoreValue().Sortable())
			doc.AddField(bluge.NewNumericField("pos", float64(d.Position)).Sortable())
			doc.AddField(bluge.NewTextField("content", d.Content).SearchTermPositions())
			doc.AddField(bluge.NewKeywordField("type", d.ContentType).StoreValue().Sortable())

			batch.Insert(doc)
		}
		if err := writer.Batch(batch); err != nil {
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
