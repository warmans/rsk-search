package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"math"
	"os"
	"path"
	"strconv"
	"time"
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
					fmt.Println("WARNING: failed to sync logger: " + err.Error())
				}
			}()

			logger.Info("Importing transcript data from...", zap.String("path", inputDir))

			dirEntries, err := os.ReadDir(inputDir)
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
				if len(episode.Transcript) == 0 {
					continue
				}

				logger.Info("Processing file...", zap.String("path", dirEntry.Name()))

				episodeDuration := getEpisodeDuration(episode)
				if episodeDuration == 0 {
					continue
				}

				wpm := calculateWordsPerSecond(episodeDuration, episode.Transcript)

				numAccurateOffsets := float64(0)
				for lineNum := range episode.Transcript {
					// calculate the missing offsets
					episode.Transcript[lineNum].Timestamp, episode.Transcript[lineNum].TimestampInferred = wpm.getSecondOffset(int64(lineNum))
					if !episode.Transcript[lineNum].TimestampInferred {
						numAccurateOffsets++
					}

					// calculate the distance to the closest non-inferred offset
					episode.Transcript[lineNum].TimestampDistance = wpm.getOffsetDistance(int64(lineNum))

					//calculate the duration of the previous line
					if lineNum >= 1 {
						episode.Transcript[lineNum-1].Duration = episode.Transcript[lineNum].Timestamp - episode.Transcript[lineNum-1].Timestamp
					}
				}
				// the final line's duration must be calculated based on the episode duration
				episode.Transcript[len(episode.Transcript)-1].Duration = episodeDuration - episode.Transcript[len(episode.Transcript)-1].Timestamp

				episode.OffsetAccuracy = int32(numAccurateOffsets / float64(len(episode.Transcript)) * 100)

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

// returns the offset + true if it was inferred (false if it is a real/accurate offset)
func (w speechVelocity) getSecondOffset(lineNum int64) (time.Duration, bool) {
	if r, ok := w.rangeIndex(lineNum); ok {
		totalOffset := w.ranges[r].startTimestamp
		relativeLineNum := lineNum - w.ranges[r].firstLineNum
		for k, v := range w.ranges[r].lineDurations {
			if int64(k) < relativeLineNum {
				totalOffset += v
			}
		}
		if totalOffset == w.ranges[r].startTimestamp {
			return totalOffset, false
		}
		return totalOffset, true
	}
	return -1, true
}

// returns the distance to the nearest offset from the given line number
func (w speechVelocity) getOffsetDistance(lineNum int64) int64 {
	if r, ok := w.rangeIndex(lineNum); ok {
		distanceToPreviousOffset := lineNum - w.ranges[r].firstLineNum
		distanceToNextOffset := w.ranges[r].lastLineNum - distanceToPreviousOffset
		return min(distanceToNextOffset, distanceToPreviousOffset)
	}
	return math.MaxInt64
}

func (w speechVelocity) rangeIndex(lineNum int64) (int, bool) {
	for k, v := range w.ranges {
		if lineNum >= v.firstLineNum && lineNum <= v.lastLineNum {
			return k, true
		}
	}
	return 0, false
}

func (w speechVelocity) currentStartSecond() time.Duration {
	if len(w.ranges) == 0 {
		return 0
	}
	return w.ranges[len(w.ranges)-1].startTimestamp
}

type offsetRange struct {
	startTimestamp time.Duration
	firstLineNum   int64
	lastLineNum    int64
	duration       time.Duration
	totalChars     int64
	lineDurations  []time.Duration
}

func getEpisodeDuration(ep *models.Transcript) time.Duration {
	if durationMsStr, ok := ep.Meta["duration_ms"]; ok {
		ms, err := strconv.Atoi(durationMsStr)
		if err != nil {
			return 0
		}
		return time.Duration(ms) * time.Millisecond
	}
	return 0
}

func calculateWordsPerSecond(totalLength time.Duration, dialog []models.Dialog) speechVelocity {
	vel := speechVelocity{
		ranges: []offsetRange{
			{
				startTimestamp: dialog[0].Timestamp, //sometimes there is an offset at the start of the dialog
				firstLineNum:   0,
				lastLineNum:    0,
				duration:       0,
				lineDurations:  []time.Duration{},
			},
		},
	}
	for lineNum, line := range dialog {
		if line.Timestamp != 0 && line.Timestamp != vel.currentStartSecond() && !line.TimestampInferred {

			// finalize current range
			vel.ranges[len(vel.ranges)-1].lastLineNum = int64(lineNum) - 1
			vel.ranges[len(vel.ranges)-1].duration = line.Timestamp - vel.currentStartSecond()

			// start a new range
			vel.ranges = append(
				vel.ranges,
				offsetRange{
					startTimestamp: line.Timestamp,
					firstLineNum:   int64(lineNum),
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
			vel.ranges[len(vel.ranges)-1].duration = totalLength - vel.ranges[len(vel.ranges)-1].startTimestamp
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
		lineDuration := float64(vel.ranges[rangeIdx].duration.Milliseconds()) * lineProportion
		vel.ranges[rangeIdx].lineDurations = append(vel.ranges[rangeIdx].lineDurations, time.Duration(lineDuration)*time.Millisecond)
	}

	return vel
}
