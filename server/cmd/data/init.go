package data

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func InitCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create thee base data files with the basic required metadata",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()
			for date, name := range meta.EpisodeDates() {
				if err := initEpisodeFile(logger, cfg, date, name); err != nil {
					return fmt.Errorf("failed to init file for date: %s name: %s: %w", date, name, err)
				}
			}
			return nil
		},
	}

	return cmd
}

func initEpisodeFile(logger *zap.Logger, cfg dataConfig, dateStr string, name string) error {
	d, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return err
	}
	publication, series, episode, err := models.ParseEpID(name)
	if err != nil {
		return err
	}
	ep := &models.Transcript{
		Publication: publication,
		Series:      series,
		Episode:     episode,
		ReleaseDate: d,
		Version:     "0.0.0",
		Transcript:  []models.Dialog{},
		Meta: map[models.MetadataType]string{
			models.CoverArtURL: "/assets/cover/default.jpg",
		},
	}

	filePath := data.EpisodeFileName(cfg.dataDir, ep)
	if ok, err := util.FileExists(filePath); ok || err != nil {
		logger.Info("Exists...", zap.String("path", filePath))
		return err
	}

	logger.Info("Creating...", zap.String("episode", name))

	if cfg.audioDir != "" {
		durationMs, err := getAudioDurationMs(
			path.Join(
				cfg.audioDir,
				fmt.Sprintf("%s-%s.mp3", publication, models.FormatStandardEpisodeName(series, episode)),
			),
		)
		if err != nil {
			logger.Warn("failed to get episode duration", zap.String("name", name), zap.Error(err))
		} else {
			ep.Meta[models.MetadataTypeDurationMs] = fmt.Sprintf("%d", durationMs)
		}
	} else {
		logger.Warn("cannot get episode duration as file path not specified", zap.String("name", name))
	}

	return data.SaveEpisodeToFile(cfg.dataDir, ep)
}

func getAudioDurationMs(audioFilePath string) (int64, error) {
	cmd := exec.Command("mp3info", "-p", "%S", audioFilePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to shell out to mp3info (is it installed?) for file: %s (raw: %s)", audioFilePath, string(out))
	}
	intVal, err := strconv.Atoi(string(out))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to convert mp3info output to int: %s", string(out))
	}
	return int64(intVal) * 1000, nil
}
