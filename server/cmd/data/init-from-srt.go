package data

import (
	"encoding/json"
	"fmt"
	"github.com/konifar/go-srt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var htmlTag = regexp.MustCompile(`<[^<>]+>`)
var lineWithActor = regexp.MustCompile(`[a-zA-Z0-9]+:.+`)

// InitFromSrtCmd
// Examples:
// for i in $(seq 1 1); do sed -i '1s/^\xEF\xBB\xBF//' ~/audio-src/source/scrimpton-audio/video/extras-S01E0${i}.srt; ./bin/rsk-search data init-from-srt --srt-path ~/audio-src/source/scrimpton-audio/video/extras-S01E0${i}.srt -p extras -s 1 -e ${i} -m ~/audio-src/source/scrimpton-audio/video/extras-S01E0${i}.mp4; done

func InitFromSrtCmd() *cobra.Command {

	var srtPath string
	var mediaFilePath string
	var publication string
	var series int32
	var episode int32

	cmd := &cobra.Command{
		Use:   "init-from-srt",
		Short: "Generate metadata files from an SRT.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			f, err := os.Open(srtPath)
			if err != nil {
				return err
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					logger.Error("failed to close file", zap.Error(err))
				}
			}(f)

			logger.Info(fmt.Sprintf("processing %s", srtPath))
			return initEpisodeFileFromSRT(logger, mediaFilePath, publication, series, episode, cfg.dataDir, f)
		},
	}

	cmd.Flags().StringVarP(&srtPath, "srt-path", "", "", "Path to SRT file")
	cmd.Flags().StringVarP(&mediaFilePath, "media-file-name", "m", "", "name of the media file")
	cmd.Flags().StringVarP(&publication, "publication", "p", "other", "Publication to give episodes")
	cmd.Flags().Int32VarP(&series, "series", "s", 1, "use this as the series number in meta/renamed file")
	cmd.Flags().Int32VarP(&episode, "episode", "e", 1, "episode number")
	return cmd
}

func initEpisodeFileFromSRT(
	logger *zap.Logger,
	mediaFileName string,
	publication string,
	series int32,
	episode int32,
	dataDir string,
	srtFile io.Reader,
) error {

	info, err := getffProbeMedia(path.Join(cfg.videoDir, mediaFileName))
	if err != nil {
		return err
	}
	durationSeconds, err := strconv.ParseFloat(info.Streams[0].Duration, 64)
	if err != nil {
		return err
	}

	ep := &models.Transcript{
		MediaType:     models.MediaTypeVideo,
		MediaFileName: path.Base(mediaFileName),
		Publication:   publication,
		Series:        series,
		Episode:       episode,
		Version:       "0.0.0",
		Transcript:    []models.Dialog{},
		Locked:        false,
		Meta: map[models.MetadataType]string{
			models.CoverArtURL:            "/assets/cover/default.jpg",
			models.MetadataTypeDurationMs: fmt.Sprintf("%d", int64(durationSeconds*1000)),
		},
		AudioQuality: models.AudioQualityGood,
	}

	scanner := gosrt.NewScanner(srtFile)
	scanned := 0
	for scanner.Scan() {
		sub := scanner.Subtitle()
		actor := "Unknown"
		if lineWithActor.MatchString(sub.Text) {
			lineParts := strings.SplitN(sub.Text, ":", 2)
			actor = lineParts[0]
		}
		ep.Transcript = append(ep.Transcript, models.Dialog{
			ID:        fmt.Sprintf("ep-%s-%d", ep.ShortID(), sub.Number),
			Type:      models.DialogTypeChat,
			Position:  int64(sub.Number),
			Actor:     actor,
			Content:   strings.TrimSpace(strings.TrimPrefix(stripTags(sub.Text), fmt.Sprintf("%s:", actor))),
			Timestamp: sub.Start,
			Duration:  sub.End - sub.Start,
		})
		scanned++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner failed after %d lines: %w", scanned, err)
	}

	filePath := data.EpisodeFileName(dataDir, ep)

	if ok, err := util.FileExists(filePath); ok || err != nil {
		logger.Info("Exists...", zap.String("path", filePath))
		if ok && err == nil {
			err = fmt.Errorf("file already exists")
		}
		return err
	}

	logger.Info("Creating...", zap.String("episode", ep.ShortID()))

	return data.SaveEpisodeToFile(dataDir, ep)
}

func stripTags(s string) string {
	return htmlTag.ReplaceAllString(s, "")
}

type ffprobeInfo struct {
	Streams []struct {
		Index              int    `json:"index"`
		CodecName          string `json:"codec_name"`
		CodecLongName      string `json:"codec_long_name"`
		Profile            string `json:"profile"`
		CodecType          string `json:"codec_type"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTag           string `json:"codec_tag"`
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		CodedWidth         int    `json:"coded_width"`
		CodedHeight        int    `json:"coded_height"`
		ClosedCaptions     int    `json:"closed_captions"`
		HasBFrames         int    `json:"has_b_frames"`
		SampleAspectRatio  string `json:"sample_aspect_ratio"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
		PixFmt             string `json:"pix_fmt"`
		Level              int    `json:"level"`
		ColorRange         string `json:"color_range"`
		ColorSpace         string `json:"color_space"`
		ColorTransfer      string `json:"color_transfer"`
		ColorPrimaries     string `json:"color_primaries"`
		ChromaLocation     string `json:"chroma_location"`
		Refs               int    `json:"refs"`
		IsAvc              string `json:"is_avc"`
		NalLengthSize      string `json:"nal_length_size"`
		RFrameRate         string `json:"r_frame_rate"`
		AvgFrameRate       string `json:"avg_frame_rate"`
		TimeBase           string `json:"time_base"`
		StartPts           int    `json:"start_pts"`
		StartTime          string `json:"start_time"`
		DurationTs         int    `json:"duration_ts"`
		Duration           string `json:"duration"`
		BitRate            string `json:"bit_rate"`
		BitsPerRawSample   string `json:"bits_per_raw_sample"`
		NbFrames           string `json:"nb_frames"`
		Disposition        struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
		} `json:"disposition"`
		Tags struct {
			CreationTime time.Time `json:"creation_time"`
			Language     string    `json:"language"`
			HandlerName  string    `json:"handler_name"`
			VendorId     string    `json:"vendor_id"`
		} `json:"tags"`
	} `json:"streams"`
}

func getffProbeMedia(videoFilePath string) (*ffprobeInfo, error) {
	// ffprobe -v quiet -show_streams -select_streams v:0 -of json
	cmd := exec.Command("ffprobe", "-v", "quiet", "-show_streams", "-select_streams", "v:0", "-of", "json", videoFilePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to shell out to ffprobe (is it installed?) for file: %s (raw: %s)", videoFilePath, string(out))
	}

	probeInfo := &ffprobeInfo{}
	if err := json.Unmarshal(out, probeInfo); err != nil {
		return nil, err
	}
	return probeInfo, nil
}
