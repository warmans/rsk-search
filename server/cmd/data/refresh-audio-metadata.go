package data

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type audioMeta struct {
	bitrateKbps float64
}

// RefreshAudioMetadataCmd
// e.g. ./bin/rsk-search data refresh-audio-metadata -a ${MEDIA_BASE_PATH}/episode
func RefreshAudioMetadataCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "refresh-audio-metadata",
		Short: "update episode files with metadata from audio files",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()
			if cfg.audioDir == "" {
				logger.Fatal("Audio dir not specified")
			}

			logger.Info("Importing transcript data from...", zap.String("path", cfg.dataDir))

			audioFileMeta, err := fetchAudioFileMeta(cfg.audioDir)
			if err != nil {
				logger.Fatal("Failed to index audio files")
			}
			dirEntries, err := os.ReadDir(cfg.dataDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {
				if dirEntry.IsDir() || !strings.HasSuffix(dirEntry.Name(), ".json") {
					continue
				}
				episode := &models.Transcript{}
				if err := util.WithReadJSONFileDecoder(path.Join(cfg.dataDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				if meta, ok := audioFileMeta[episode.ShortID()]; ok {
					episode.Meta[models.MetadataTypeBitrateKbps] = fmt.Sprintf("%0.2f", meta.bitrateKbps)
					if err := data.ReplaceEpisodeFile(cfg.dataDir, episode); err != nil {
						return err
					}
				} else {
					logger.Info("Failed to find bitrate for file", zap.String("id", episode.ShortID()))
				}
			}
			return nil
		},
	}

	return cmd
}

func fetchAudioFileMeta(audioFileDir string) (map[string]audioMeta, error) {
	dirEntries, err := os.ReadDir(audioFileDir)
	if err != nil {
		return nil, err
	}
	index := make(map[string]audioMeta)
	for _, e := range dirEntries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".mp3") {
			continue
		}
		bitRateKbps, err := getBitrate(path.Join(audioFileDir, e.Name()))
		if err != nil {
			return nil, err
		}
		index[strings.TrimSuffix(e.Name(), ".mp3")] = audioMeta{
			bitrateKbps: bitRateKbps,
		}
	}
	return index, nil
}

func getBitrate(audioFilePath string) (float64, error) {
	cmd := exec.Command("mp3info", "-r", "a", "-p", "%r", audioFilePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to shell out to mp3info (is it installed?) for file: %s (raw: %s)", audioFilePath, string(out))
	}
	asFloat, err := strconv.ParseFloat(string(out), 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to convert mp3info output to float: %s", string(out))
	}
	return asFloat, nil
}
