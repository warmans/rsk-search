package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/spotify"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func ImportSpotifyData() *cobra.Command {

	var tinPotRadioData string
	var spotifyDataPath string
	var forceCacheRefresh bool

	cmd := &cobra.Command{
		Use:   "spotifize",
		Short: "Import data from spotify",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
			spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

			var spotifyToken string
			if spotifyClientID == "" || spotifyClientSecret == "" {
				logger.Warn("No spotify client ID or secret configured. Local cache will be used instead.")
			} else {
				resp, err := http.Post(
					"https://accounts.spotify.com/api/token",
					"application/x-www-form-urlencoded",
					strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", spotifyClientID, spotifyClientSecret)),
				)
				if err != nil {
					logger.Error("Failed to get spotify token", zap.Error(err))
					return err
				}
				defer resp.Body.Close()

				token := accessToken{}
				if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
					return err
				}
				spotifyToken = token.AccessToken
				logger.Info("Spotify token was created...")
			}

			if err := addTinPotRadioLinks(tinPotRadioData, logger.With(zap.String("stage", "tinpotradio"))); err != nil {
				fmt.Println("Failed to add tinpot radio links (try running script first): ", err.Error())
			}

			return addSongMeta(logger.With(zap.String("stage", "songs")), spotifyToken, spotifyDataPath, forceCacheRefresh)
		},
	}

	cmd.Flags().StringVarP(&tinPotRadioData, "tinpotradio-episodes", "i", "./script/tinpotradio/raw/xfm-spotify-meta.json", "Path to tinpot radio data (episode links)")
	cmd.Flags().StringVarP(&spotifyDataPath, "spotify-data", "s", "./pkg/meta/data/songs.json", "Path to tinpot radio data (episode links)")
	cmd.Flags().BoolVarP(&forceCacheRefresh, "force-cache-refresh", "f", false, "if true cache will be ignored and all data refetched from spotify")

	return cmd
}

func addTinPotRadioLinks(tinPotRadioData string, logger *zap.Logger) error {

	logger.Info("Importing tinpotradio episodes from...", zap.String("path", tinPotRadioData))

	f, err := os.Open(tinPotRadioData)
	if err != nil {
		return err
	}
	defer f.Close()

	episodes := []spotify.Episode{}
	if err := json.NewDecoder(f).Decode(&episodes); err != nil {
		return err
	}

	for _, d := range episodes {

		lg := logger.With(zap.String("name", d.Name), zap.String("publication", meta.PublicationXFM), zap.String("date", d.ReleaseDate))

		ep, err := data.LoadEpisodeByName(cfg.dataDir, meta.PublicationXFM, d.Name)
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

		if err := data.ReplaceEpisodeFile(cfg.dataDir, ep); err != nil {
			return err
		}
		lg.Info("ok!")
	}
	return nil
}

func addSongMeta(logger *zap.Logger, token string, metadataPath string, forceCacheRefresh bool) error {

	search := spotify.NewSearch(token)

	f, err := os.OpenFile(metadataPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	songCache := meta.Songs{}
	if err := json.NewDecoder(f).Decode(&songCache); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	for _, shortId := range meta.EpisodeDates() {

		lg := logger.With(zap.String("name", shortId))

		ep, err := data.LoadEpisodeByShortID(cfg.dataDir, shortId)
		if err != nil {
			return err
		}
		if ep == nil {
			lg.Info("No episode was initialized for this spotify result")
			continue
		}

		for k, v := range ep.Transcript {
			if v.Type == models.DialogTypeSong {
				searchTerm := util.StripNonAlphanumeric(strings.TrimSpace(v.Content))
				if len(searchTerm) < 3 || strings.ToLower(searchTerm) == "unknown" {
					continue
				}

				var track *spotify.Track

				cachedId, ok := songCache.FindKeyByTerm(searchTerm)
				if !ok || forceCacheRefresh {
					if token == "" {
						lg.Warn("No spotify token given. Using cache only...", zap.String("term", searchTerm))
						continue
					}
					lg.Info("Querying spotify...", zap.String("term", searchTerm))
					track, err = search.FindTrack(searchTerm)
					if err != nil || track == nil {
						lg.Warn("failed to find track", zap.Error(err), zap.String("term", searchTerm))
						continue
					}
					if songCache.Songs == nil {
						songCache.Songs = make(map[string]*meta.Song)
					}
					songCache.Songs[track.TrackURI] = &meta.Song{
						Terms:      []string{searchTerm},
						EpisodeIDs: []string{shortId},
						Track:      track,
					}
				} else {
					lg.Info("Cached...", zap.String("term", searchTerm))

					track = songCache.Songs[cachedId].Track
					var found bool
					for _, epID := range songCache.Songs[cachedId].EpisodeIDs {
						if epID == shortId {
							found = true
						}
					}
					if !found {
						songCache.Songs[cachedId].EpisodeIDs = append(songCache.Songs[cachedId].EpisodeIDs, shortId)
					}

					found = false
					for _, term := range songCache.Songs[cachedId].Terms {
						if term == searchTerm {
							found = true
						}
					}
					if !found {
						songCache.Songs[cachedId].Terms = append(songCache.Songs[cachedId].Terms, searchTerm)
					}
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

		if err := data.ReplaceEpisodeFile(cfg.dataDir, ep); err != nil {
			return err
		}
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(songCache)
}
