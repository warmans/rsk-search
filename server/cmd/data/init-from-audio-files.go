package data

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
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
	"strconv"
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
}

func InitFromAudioFilesCmd() *cobra.Command {

	var publication string
	var metadataStrategy string

	cmd := &cobra.Command{
		Use:   "init-from-audio",
		Short: "Generate metadata files from audio files (by their name).",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			audioFiles := []audioFile{}

			entries, err := os.ReadDir(cfg.audioDir)
			if err != nil {
				return err
			}
			for _, v := range entries {
				if v.IsDir() {
					continue
				}
				filePath := path.Join(cfg.audioDir, v.Name())

				if metadataStrategy == "filename" {
					dateStr, name, year, err := parseFileName(logger, v.Name())
					if err != nil {
						logger.Warn(fmt.Sprintf("Failed to parse %s", v.Name()))
						continue
					}
					audioFiles = append(audioFiles, audioFile{
						path:        filePath,
						name:        name,
						date:        dateStr,
						year:        year,
						publication: publication,
					})
				} else {
					meta, err := parseMetadata(logger, filePath, publication)
					if err != nil {
						return err
					}
					audioFiles = append(audioFiles, meta)
				}
			}

			sort.Slice(audioFiles, func(i, j int) bool {
				return audioFiles[i].year < audioFiles[i].year
			})

			renamedFileDir := path.Join(cfg.audioDir, "renamed")
			for k, f := range audioFiles {
				ep, err := initEpisodeFileFromAudio(logger, f, int32(k)+1, cfg.dataDir)
				if err != nil {
					return fmt.Errorf("failed to init file for file %s date: %s name: %s: %w", f.path, f.date, f.name, err)
				}
				if _, err := exec.Command("cp", f.path, path.Join(renamedFileDir, fmt.Sprintf("%s.mp3", ep.ShortID()))).CombinedOutput(); err != nil {
					return err
				}

			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&publication, "publication", "p", "other", "Publication to give episodes")
	cmd.Flags().StringVarP(&metadataStrategy, "meta-strategy", "m", "id3", "id3 or filename")

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

	spew.Dump(rawTags)

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

// 2005 - Extras - Steve Interviewed by Simon Amstell on XFM 2005-08-13.mp3
func parseFileName(logger *zap.Logger, fileName string) (*time.Time, string, int, error) {
	fileName = strings.TrimSuffix(fileName, ".mp3")
	segments := strings.Split(fileName, " ")

	// last segment should be a date in the format YYY-MM-DD or YYY-MM
	dateParts := strings.Split(segments[len(segments)-1], "-")
	if len(dateParts) == 2 {
		dateParts = append(dateParts, "01")
	}
	dateStr := fmt.Sprintf("%sT00:00:00Z", strings.Join(dateParts, "-"))
	name := strings.Join(segments[:len(segments)-1], " ")

	ts, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		logger.Warn(fmt.Sprintf("%s has an invalid timestamp: %s", fileName, dateStr))
	}

	year := strings.Replace(segments[0], "x", "0", -1)
	intYear, err := strconv.Atoi(year)
	if err != nil {
		return nil, "", 0, err
	}

	return &ts, name, intYear, nil
}

func initEpisodeFileFromAudio(
	logger *zap.Logger,
	f audioFile,
	episode int32,
	dataDir string,
) (*models.Transcript, error) {

	ep := &models.Transcript{
		Publication: f.publication,
		Series:      1,
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
	}

	return ep, data.SaveEpisodeToFile(dataDir, ep)
}
