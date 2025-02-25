package data

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	ffmpeg_go "github.com/warmans/ffmpeg-go/v2"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
)

// ExtractVideoImages
// create sprites
func ExtractVideoImages() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "extract-video-images",
		Short: "extract an image for every line in a video transcript",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()
			if cfg.videoDir == "" {
				logger.Fatal("Video dir not specified")
			}

			logger.Info("Importing transcript data from...", zap.String("path", cfg.dataDir))

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
				if episode.Media.VideoFileName == "" {
					continue
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				for _, t := range episode.Transcript {
					if t.Position%5 == 0 {
						if err := dumpImageIfNotExists(cfg.imageDir, t, path.Join(cfg.videoDir, episode.MediaFileName)); err != nil {
							//todo: don't fail
							logger.Fatal("failed to dump image", zap.String("id", t.ID), zap.Error(err))
						}
					}

				}
			}
			return nil
		},
	}

	return cmd
}

func dumpImageIfNotExists(imagesDir string, d models.Dialog, videoPath string) error {

	outputPath := path.Join(imagesDir, fmt.Sprintf("%s.png", d.ID))

	//skip existing files
	_, err := os.Stat(outputPath)
	if err == nil {
		return nil
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return ffmpeg_go.
		Input(videoPath,
			ffmpeg_go.KwArgs{
				"ss": fmt.Sprintf("%0.2f", d.Timestamp.Seconds()),
			}).
		Output(outputPath,
			ffmpeg_go.KwArgs{
				"frames:v":       "1",
				"format":         "apng",
				"filter_complex": "[0:v]scale=164:-1",
			},
		).WithErrorOutput(os.Stderr).Run()

}
