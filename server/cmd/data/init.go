package data

import (
	"fmt"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"go.uber.org/zap"
	"time"
)

func InitCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create thee base data files with the basic required metadata",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			for date, name := range meta.EpisodeDates() {
				if err := initEpisodeFile(cfg.dataDir, date, name); err != nil {
					return fmt.Errorf("failed to init file for date: %s name: %s: %w", date, name, err)
				}
			}
			return nil
		},
	}

	return cmd
}

func initEpisodeFile(outputDir string, dateStr string, name string) error {
	d, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return err
	}

	series, episode, err := models.ParseStandardEpisodeName(name)
	if err != nil {
		return err
	}

	ep := &models.Transcript{
		Publication: meta.PublicationXFM,
		Series:      series,
		Episode:     episode,
		ReleaseDate: d,
	}
	return data.SaveEpisodeToFile(outputDir, ep)
}
