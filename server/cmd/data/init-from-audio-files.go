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
	"regexp"
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
	var dryRun bool

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

				logger.Info(fmt.Sprintf("processing %s", filePath))

				meta, err := parseFileName(logger, filePath, v.Name(), publication)
				if err != nil {
					logger.Warn(fmt.Sprintf("Failed to parse filename %s, fall back to id3", v.Name()))
					meta, err = parseMetadata(logger, filePath, publication)
					if err != nil {
						logger.Warn(fmt.Sprintf("Failed to parse id3 of %s, giving up", v.Name()))
						continue
					}
					audioFiles = append(audioFiles, meta)
				} else {
					audioFiles = append(audioFiles, meta)
				}
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

			renamedFileDir := path.Join(cfg.audioDir, "renamed")
			if _, err := exec.Command("rm", "-rf", fmt.Sprintf("%s/*", renamedFileDir)).CombinedOutput(); err != nil {
				return err
			}
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
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "x", false, "don't write any files ")

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

// 2005 - Extras - Steve Interviewed by Simon Amstell on XFM 2005-08-13.mp3
func parseFileName(logger *zap.Logger, filePath string, fileName string, publication string) (audioFile, error) {

	var year string
	var name string
	var date string

	fileName = strings.TrimSpace(fileName)
	fullPattern := regexp.MustCompile("^([0-9x]+)?[\\s\\-]*(.+)([0-9]{4}-[0-9]{2}-[0-9]{2}).+$")
	matches := fullPattern.FindAllStringSubmatch(fileName, -1)

	if len(matches) == 0 {
		partialPattern := regexp.MustCompile("^([0-9x]+)?[\\s\\-]*(.+)$")
		matches = partialPattern.FindAllStringSubmatch(fileName, -1)

		year = matches[0][1]
		name = matches[0][2]
		date = ""

		if len(matches) == 0 {
			return audioFile{}, fmt.Errorf("name does not match")
		}
	} else {
		year = matches[0][1]
		name = matches[0][2]
		date = matches[0][3]
	}

	var ts time.Time
	if date != "" {
		dateStr := fmt.Sprintf("%sT00:00:00Z", date)
		var err error
		ts, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			logger.Warn(fmt.Sprintf("%s has an invalid timestamp: %s", fileName, dateStr))
		}
	}

	intYear, err := strconv.Atoi(strings.Replace(year, "x", "0", -1))
	if err != nil {
		logger.Warn(fmt.Sprintf("%s has an invalid year: %s", fileName, date))
	}

	return audioFile{
		path:        filePath,
		name:        strings.TrimSuffix(strings.TrimSpace(name), ".mp3"),
		date:        timePointer(ts),
		year:        intYear,
		publication: publication,
	}, nil
}

func timePointer(ts time.Time) *time.Time {
	if ts.IsZero() {
		return nil
	}
	return &ts
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
