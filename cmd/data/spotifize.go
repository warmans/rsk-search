package data

import (
	"encoding/json"
	"fmt"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
)

func ImportSpotifyData() *cobra.Command {

	var inputFile string

	cmd := &cobra.Command{
		Use:   "spotifize",
		Short: "Import data from spotify",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			logger.Info("Importing spotify data from...", zap.String("path", inputFile))

			f, err := os.Open(inputFile)
			if err != nil {
				return err
			}
			defer f.Close()

			data := []models.SpotifyItem{}
			if err := json.NewDecoder(f).Decode(&data); err != nil {
				return err
			}

			for _, d := range data {
				lg := logger.With(zap.String("name", d.Name), zap.String("publication", meta.PublicationXFM), zap.String("date", d.ReleaseDate))

				ep, err := util.LoadEpisode(cfg.dataDir, meta.PublicationXFM, d.Name)
				if err != nil {
					return err
				}
				if ep == nil {
					lg.Info("No episode was initialized for this spotify result")
					continue
				}

				if ep.Meta == nil {
					ep.Meta = map[models.MetadataType]string{}
				}
				ep.Meta[models.MetadataTypeDurationMs] = fmt.Sprint(d.DurationMs)
				ep.Meta[models.MetadataTypeSpotifyURI] = d.URI
				ep.Meta[models.MetadataTypeSpotifyPreviewURL] = d.AudioPreviewURL

				if err := util.SaveEpisodeToFile(cfg.dataDir, ep); err != nil {
					return err
				}
				lg.Info("ok!")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "./script/tinpotradio/raw/xfm-spotify-meta.json", "Path to raw scraped file")

	return cmd
}
