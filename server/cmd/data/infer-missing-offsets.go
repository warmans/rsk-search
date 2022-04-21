package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io/ioutil"
	"math"
	"path"
	"strconv"
)

func InferMissingOffsetsCmd() *cobra.Command {

	var inputDir string
	var singleEpisode string

	cmd := &cobra.Command{
		Use:   "infer-missing-offsets",
		Short: "adds time offsets to dialog based on approximate words per min spoken",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger, _ := zap.NewProduction()
			defer func() {
				if err := logger.Sync(); err != nil {
					fmt.Println("WARNING: failed to sync logger: "+err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			dirEntries, err := ioutil.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, dirEntry := range dirEntries {

				if dirEntry.IsDir() {
					continue
				}

				episode := &models.Transcript{}
				if err := util.WithReadJSONFileDecoder(path.Join(inputDir, dirEntry.Name()), func(dec *json.Decoder) error {
					return dec.Decode(episode)
				}); err != nil {
					return err
				}

				if singleEpisode != "" {
					if episode.ShortID() != singleEpisode {
						continue
					}
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				episodeDuration := getEpisodeDuration(episode)
				if episodeDuration == 0 {
					continue
				}

				wpm := calculateWordsPerSecond(episodeDuration, episode.Transcript)

				numAccurateOffsets := float64(0)
				for lineNum := range episode.Transcript {
					episode.Transcript[lineNum].OffsetSec, episode.Transcript[lineNum].OffsetInferred = wpm.getSecondOffset(int64(lineNum))
					if !episode.Transcript[lineNum].OffsetInferred {
						numAccurateOffsets++
					}
				}
				episode.OffsetAccuracy =  int32(numAccurateOffsets / float64(len(episode.Transcript)) * 100)

				if err := data.ReplaceEpisodeFile(inputDir, episode); err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputDir, "input-path", "i", "./var/data/episodes", "Path to raw scraped files")
	cmd.Flags().StringVarP(&singleEpisode, "single-episode", "s", "", "Only process the given episode e.g. xfm-S2E04")

	return cmd
}

type speechVelocity struct {
	ranges []offsetRange
}

// returns the offset + true if it was inferred (false if it is an real/accurate offset)
func (w speechVelocity) getSecondOffset(lineNum int64) (int64, bool) {
	if r, ok := w.rangeIndex(lineNum); ok {
		totalOffset := float64(w.ranges[r].startSecond)
		relativeLineNum := lineNum - w.ranges[r].firstLineNum
		for k, v := range w.ranges[r].lineDurations {
			if int64(k) < relativeLineNum {
				totalOffset += v
			}
		}
		if totalOffset == float64(w.ranges[r].startSecond) {
			return int64(math.Round(totalOffset)), false
		}
		return int64(math.Ceil(totalOffset)), true
	}
	return -1, true
}

func (w speechVelocity) rangeIndex(lineNum int64) (int, bool) {
	for k, v := range w.ranges {
		if lineNum >= v.firstLineNum && lineNum <= v.lastLineNum {
			return k, true
		}
	}
	return 0, false
}

func (w speechVelocity) currentStartSecond() int64 {
	if len(w.ranges) == 0 {
		return 0
	}
	return w.ranges[len(w.ranges)-1].startSecond
}

type offsetRange struct {
	startSecond     int64
	firstLineNum    int64
	lastLineNum     int64
	durationSeconds int64
	totalChars      int64
	lineDurations   []float64
}

func getEpisodeDuration(ep *models.Transcript) int64 {
	if durationMsStr, ok := ep.Meta["duration_ms"]; ok {
		ms, err := strconv.Atoi(durationMsStr)
		if err != nil {
			return 0
		}
		return int64(ms / 1000)
	}
	return 0
}

func calculateWordsPerSecond(totalLengthSeconds int64, dialog []models.Dialog) speechVelocity {
	vel := speechVelocity{
		ranges: []offsetRange{
			{
				startSecond:     dialog[0].OffsetSec, //sometimes there is an offset at the start of the dialog
				firstLineNum:    0,
				lastLineNum:     0,
				durationSeconds: 0,
				lineDurations:   []float64{},
			},
		},
	}
	for lineNum, line := range dialog {
		if line.OffsetSec != 0 && line.OffsetSec != vel.currentStartSecond() && !line.OffsetInferred {

			// finalize current range
			vel.ranges[len(vel.ranges)-1].lastLineNum = int64(lineNum) - 1
			vel.ranges[len(vel.ranges)-1].durationSeconds = line.OffsetSec - vel.currentStartSecond()

			// start a new range
			vel.ranges = append(
				vel.ranges,
				offsetRange{
					startSecond:  line.OffsetSec,
					firstLineNum: int64(lineNum),
				},
			)
		}

		// count up all the chars in the speech lines
		if line.Type != models.DialogTypeSong {
			vel.ranges[len(vel.ranges)-1].totalChars += int64(len(line.Content))
		}

		// last line should always close the range.
		if lineNum == len(dialog)-1 {
			vel.ranges[len(vel.ranges)-1].lastLineNum = int64(lineNum)
			vel.ranges[len(vel.ranges)-1].durationSeconds = totalLengthSeconds - vel.ranges[len(vel.ranges)-1].startSecond
		}
	}

	for lineNum, line := range dialog {
		rangeIdx, ok := vel.rangeIndex(int64(lineNum))
		if !ok {
			continue
		}
		if line.Type == models.DialogTypeSong {
			vel.ranges[rangeIdx].lineDurations = append(vel.ranges[rangeIdx].lineDurations, 0)
			continue
		}

		charCount := float64(len(line.Content))
		lineProportion := charCount / float64(vel.ranges[rangeIdx].totalChars)
		vel.ranges[rangeIdx].lineDurations = append(vel.ranges[rangeIdx].lineDurations, float64(vel.ranges[rangeIdx].durationSeconds)*lineProportion)
	}

	return vel
}
