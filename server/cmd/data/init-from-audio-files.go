package data

import (
	"fmt"
	"github.com/dhowden/tag"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path"
	"sort"
	"time"
)

type audioFile struct {
	path        string
	name        string
	date        *time.Time
	year        int
	publication string
	summary     string
}

func InitFromAudioFilesCmd() *cobra.Command {

	var rawAudioDir string
	var publication string
	var series int32
	var episodeOffset int32
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "init-from-audio",
		Short: "Generate metadata files from audio files (by name).",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			audioFiles := []audioFile{}

			entries, err := os.ReadDir(rawAudioDir)
			if err != nil {
				return err
			}
			for _, v := range entries {
				if v.IsDir() {
					continue
				}
				filePath := path.Join(rawAudioDir, v.Name())

				logger.Info(fmt.Sprintf("processing %s", filePath))

				meta, err := parseMetadata(logger, filePath, publication)
				if err != nil {
					logger.Warn(fmt.Sprintf("Failed to parse id3 of %s, giving up: %s", v.Name(), err.Error()))
					continue
				}
				audioFiles = append(audioFiles, meta)
			}

			sort.Slice(audioFiles, func(i, j int) bool {
				return audioFiles[i].year < audioFiles[j].year
			})

			if dryRun {
				for _, f := range audioFiles {
					logger.Info("Created file", zap.String("name", f.name), zap.Timep("date", f.date), zap.Int("year", f.year))
				}
				return nil
			}

			for k, f := range audioFiles {
				ep, err := initEpisodeFileFromAudio(logger, f, series, episodeOffset+int32(k)+1, cfg.dataDir)
				if err != nil {
					return fmt.Errorf("failed to init file for file %s date: %s name: %s: %w", f.path, f.date, f.name, err)
				}
				if _, err := exec.Command("cp", f.path, path.Join(cfg.audioDir, fmt.Sprintf("%s.mp3", ep.ShortID()))).CombinedOutput(); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&rawAudioDir, "raw-audio-path", "", "", "Path to scan")
	cmd.Flags().StringVarP(&publication, "publication", "p", "other", "Publication to give episodes")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "x", false, "don't write any files, just log")
	cmd.Flags().Int32VarP(&series, "series", "s", 1, "use this as the series number in meta/renamed file")
	cmd.Flags().Int32VarP(&episodeOffset, "episode-num-offset", "e", 0, "use this as the first episode number in meta/renamed file")
	return cmd
}

func parseMetadata(logger *zap.Logger, fileName string, publication string) (audioFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return audioFile{}, err
	}
	tags, err := tag.ReadFrom(file)
	if err != nil {
		return audioFile{}, err
	}
	rawTags := tags.Raw()

	result := audioFile{}
	result.publication = publication
	result.path = fileName
	result.name = tags.Title()
	result.summary = findDescription(rawTags)
	result.year = tags.Year()

	if dateStr := findReleaseTimestamp(rawTags); dateStr != "" {
		date, err := time.Parse(time.RFC3339, dateStr)
		if err == nil {
			result.date = &date
		} else {
			logger.Warn(fmt.Sprintf("failed to parse RELEASETIME %s", dateStr))
		}
	}
	if result.date == nil && result.year > 0 {
		startOfYear := time.Date(result.year, 0, 0, 0, 0, 0, 0, time.UTC)
		result.date = &startOfYear
	}

	return result, nil
}

func findDescription(tags map[string]interface{}) string {
	for _, key := range []string{"TDS", "TT3", "DESCRIPTION", "SUBTITLE", "PODCASTDESC"} {
		if foundTag, ok := tags[key]; ok && foundTag != nil {
			if strVal, strOk := foundTag.(string); strOk && strVal != "" {
				return strVal
			}
		}
	}
	return ""
}

func findReleaseTimestamp(tags map[string]interface{}) string {
	for _, key := range []string{"TDEN", "TDRL", "TDR"} {
		if foundTag, ok := tags[key]; ok && foundTag != nil {
			if strVal, strOk := foundTag.(string); strOk {
				return strVal
			}
		}
	}
	return ""
}

func initEpisodeFileFromAudio(
	logger *zap.Logger,
	f audioFile,
	series int32,
	episode int32,
	dataDir string,
) (*models.Transcript, error) {

	ep := &models.Transcript{
		Publication: f.publication,
		Series:      series,
		Episode:     episode,
		ReleaseDate: f.date,
		Name:        f.name,
		Version:     "0.0.0",
		Summary:     f.summary,
		Transcript:  []models.Dialog{},
		Locked:      true,
		Meta: map[models.MetadataType]string{
			models.CoverArtURL: "/assets/cover/default.jpg",
		},
	}

	filePath := data.EpisodeFileName(dataDir, ep)
	if ok, err := util.FileExists(filePath); ok || err != nil {
		logger.Info("Exists...", zap.String("path", filePath))
		if ok && err == nil {
			err = fmt.Errorf("file already exists")
		}
		return nil, err
	}

	logger.Info("Creating...", zap.String("episode", f.name))
	durationMs, err := getAudioDurationMs(f.path)
	if err != nil {
		logger.Warn("failed to get episode duration", zap.String("name", f.name), zap.Error(err))
	} else {
		ep.Meta[models.MetadataTypeDurationMs] = fmt.Sprintf("%d", durationMs)
		ep.Media.AudioDurationMs = durationMs
	}

	return ep, data.SaveEpisodeToFile(dataDir, ep)
}
