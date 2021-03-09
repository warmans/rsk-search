package data

import (
	"encoding/json"
	"fmt"
	_ "github.com/blevesearch/bleve/v2/config"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/spotify"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"strings"
)

func ImportSpotifyData() *cobra.Command {

	var tinPotRadioData string
	var spotifyToken string

	cmd := &cobra.Command{
		Use:   "spotifize",
		Short: "Import data from spotify",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any

			spotifyToken = os.Getenv("SPOTIFY_TOKEN")
			if spotifyToken == "" {
				return fmt.Errorf("token cannot be empty")
			}

			if err := addTinPotRadioLinks(tinPotRadioData, logger.With(zap.String("stage", "tinpotradio"))); err != nil {
				return err
			}

			return addSongMeta(logger.With(zap.String("stage", "songs")), spotifyToken)
		},
	}

	cmd.Flags().StringVarP(&tinPotRadioData, "tinpotradio-episodes", "i", "./script/tinpotradio/raw/xfm-spotify-meta.json", "Path to tinpot radio data (episode links)")

	return cmd
}

func addTinPotRadioLinks(tinPotRadioData string, logger *zap.Logger) error {

	logger.Info("Importing tinpotradio data from...", zap.String("path", tinPotRadioData))

	f, err := os.Open(tinPotRadioData)
	if err != nil {
		return err
	}
	defer f.Close()

	data := []spotify.Episode{}
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
			ep.Meta = *new(models.Metadata)
		}
		ep.Meta[models.MetadataTypeDurationMs] = fmt.Sprint(d.DurationMs)
		ep.Meta[models.MetadataTypeSpotifyURI] = d.URI
		ep.Meta[models.MetadataTypeSpotifyPreviewURL] = d.AudioPreviewURL

		if err := util.ReplaceEpisodeFile(cfg.dataDir, ep); err != nil {
			return err
		}
		lg.Info("ok!")
	}
	return nil
}

func addSongMeta(logger *zap.Logger, token string) error {

	search := spotify.NewSearch(token)

	for _, name := range meta.XfmEpisodeNames() {

		lg := logger.With(zap.String("name", name))

		ep, err := util.LoadEpisode(cfg.dataDir, meta.PublicationXFM, name)
		if err != nil {
			return err
		}
		if ep == nil {
			lg.Info("No episode was initialized for this spotify result")
			continue
		}

		for k, v := range ep.Transcript {
			if v.Type == models.DialogTypeSong {
				searchTerm := strings.TrimSpace(v.Content)
				lg.Info("Locating song", zap.String("term", searchTerm))

				track, err := search.FindTrack(searchTerm)
				if err != nil || track == nil {
					lg.Warn("failed to find track", zap.Error(err))
					continue
				}

				if v.Meta == nil {
					ep.Transcript[k].Meta = models.Metadata{}
				}
				ep.Transcript[k].Meta[models.MetadataSongTrack] = track.Name
				ep.Transcript[k].Meta[models.MetadataTypeSpotifyURI] = track.TrackURI

				if len(track.Artists) == 1 {
					ep.Transcript[k].Meta[models.MetadataSongArtist] = track.Artists[0].Name
				}
				if track.AlbumName != "" {
					ep.Transcript[k].Meta[models.MetadataSongAlbum] = track.AlbumName
				}
			}
		}

		if err := util.ReplaceEpisodeFile(cfg.dataDir, ep); err != nil {
			return err
		}
	}

	return nil
}
