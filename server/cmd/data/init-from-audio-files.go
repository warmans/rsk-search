package data

import (
	"fmt"
	"github.com/dhowden/tag"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
	"time"
)

type audioFile struct {
	path        string
	name        string
	date        *time.Time
	year        int
	publication string
	summary     string
	series      int32
	episode     int32
}

// InitFromAudioFilesCmd e.g. ./script/meta-from-audio-files.sh ${AUDIO_DIR}/podcast-S2*
func InitFromAudioFilesCmd() *cobra.Command {

	var publicationType string
	var audioFilePath string
	var publication string
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

			if audioFilePath == "" {
				return fmt.Errorf("audio file path not specified")
			}

			logger.Info(fmt.Sprintf("processing %s", audioFilePath))

			meta, err := parseMetadata(logger, audioFilePath, publication)
			if err != nil {
				return err
			}

			_, err = initEpisodeFileFromAudio(logger, meta, cfg.dataDir, models.PublicationType(publicationType))
			if err != nil {
				return fmt.Errorf("failed to init metadata for file %s date: %s name: %s: %w", meta.path, meta.date, meta.name, err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&audioFilePath, "audio-file-path", "", "", "Path to scan")
	cmd.Flags().StringVarP(&publication, "publication", "p", "", "Publication to give episodes")
	cmd.Flags().StringVarP(&publicationType, "publication-type", "t", "podcast", "Publication type to give episodes")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "x", false, "don't write any files, just log")
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

	// try and get the publication from the file name e.g. podcast-S1E01.mp3
	filenameParts := strings.Split(path.Base(fileName), "-")
	if publication == "" && len(filenameParts) == 2 {
		publication = filenameParts[0]
	}

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

	parsedSeries, parsedEpisode, err := models.ExtractSeriesAndEpisode(strings.TrimSuffix(fileName, ".mp3"))
	if err != nil {
		return result, err
	}
	result.series = parsedSeries
	result.episode = parsedEpisode

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
	dataDir string,
	publicationType models.PublicationType,
) (*models.Transcript, error) {

	ep := &models.Transcript{
		PublicationType: publicationType,
		Publication:     f.publication,
		Series:          f.series,
		Episode:         f.episode,
		ReleaseDate:     f.date,
		Name:            f.name,
		Version:         "0.0.0",
		Summary:         f.summary,
		Transcript:      []models.Dialog{},
		Locked:          false,
		Meta: map[models.MetadataType]string{
			models.CoverArtURL: "/assets/cover/default.jpg",
		},
		Media: models.Media{
			AudioFileName: path.Base(f.path),
		},
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
