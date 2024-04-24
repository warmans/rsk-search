package data

import (
	"encoding/json"
	"fmt"
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/bbalet/stopwords"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/assemblyai"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"math"
	"path"
	"regexp"
	"strings"
	"time"
)

var punctuation = regexp.MustCompile("[^a-zA-Z0-9\\s]+")

func MergeTimestampsAAICommand() *cobra.Command {

	var timestampSourceFile string
	var targetTranscriptName string
	var outputPath string
	var verbose bool
	var replace bool

	cmd := &cobra.Command{
		Use:   "merge-timestamps-aai",
		Short: "create a machine transcription using Assembly AI API",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetPath := path.Join(cfg.dataDir, targetTranscriptName)
			target, err := data.LoadEpisodePath(targetPath)
			if err != nil {
				return fmt.Errorf("failed to load target transcript at path %s: %w", targetPath, err)
			}

			timestampSource := assemblyai.TranscriptionStatusResponse{}
			if err := util.WithReadJSONFileDecoder(timestampSourceFile, func(dec *json.Decoder) error {
				return dec.Decode(&timestampSource)

			}); err != nil {
				return fmt.Errorf("failed to read timestamp source %s: %w", timestampSourceFile, err)
			}

			if outputPath == "" {
				outputPath = path.Join(cfg.dataDir, targetTranscriptName)
			}

			target.Transcript = mergeTimestampsTo(target.Transcript, assemblyAiToDialog(timestampSource.Utterances), verbose)
			if replace {
				err = data.ReplaceEpisodeFile(cfg.dataDir, target)
			} else {
				err = util.WithReplaceJSONFileEncoder(outputPath, func(enc *json.Encoder) error {
					return enc.Encode(target)
				})
			}

			if err != nil {
				return fmt.Errorf("failed to write result: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&timestampSourceFile, "timestamp-source", "s", "", "Source of timestamp data from assembly-ai")
	cmd.Flags().StringVarP(&targetTranscriptName, "target-transcript", "t", "", "Target transcript")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output result to (defaults to target)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show possible matches")
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "replace source file")

	return cmd
}

func mergeTimestampsTo(target []models.Dialog, compare []models.Dialog, verbose bool) []models.Dialog {

	// clear all non-chat data
	transcript := []models.Dialog{}
	for _, v := range target {
		if v.Type == models.DialogTypeChat && v.Actor != "" {
			transcript = append(transcript, v)
		}
	}

	lastMatchedTimestamp := time.Duration(0)
	numMatched := 0
	for targetPos, targetLine := range transcript {
		targetText := cleanString(targetLine.Content)
		numTargetWords := len(strings.Split(targetText, " "))

		fmt.Printf("TARGET %d: %s (%s)\n", targetPos, targetLine.Content, targetText)
		for comparePos := 0; comparePos < len(compare); comparePos++ {
			compareText := cleanString(compare[comparePos].Content)
			compareTimestamp := compare[comparePos].Timestamp
			numCompareWords := len(strings.Split(compareText, " "))

			// do not backtrack
			if compareTimestamp < lastMatchedTimestamp {
				continue
			}

			// skip obviously useless data
			if compareText == "" || targetText == "" || compareText == "yeah yeah" || numCompareWords == 1 || numTargetWords == 1 {
				continue
			}
			if verbose && distanceWithinNPcnt(targetPos, comparePos, 0.1) {
				fmt.Printf("\tCOMPARE %d: %s (%s)\n", comparePos, compare[comparePos].Content, compareText)
			}

			distance := distancePcnt(targetPos, comparePos)
			var matched bool
			var similarity float64
			if numTargetWords <= 4 {
				// compare the original text instead
				similarity = strutil.Similarity(targetLine.Content, compare[comparePos].Content, metrics.NewHamming())
				matched = similarity >= 0.60 && distance < 0.40
			} else {
				if len(compareText) > len(targetText) {
					// do a prefix match since the start of the text is the important bit for the timestamp
					compareText = compareText[:len(targetText)]
					similarity = strutil.Similarity(targetText, compareText, metrics.NewHamming())
					matched = similarity >= 0.70 && distance < 0.35
				} else {
					// to do a full match
					similarity = strutil.Similarity(targetText, compareText, metrics.NewHamming())
					matched = similarity >= 0.40 && distance < 0.35
				}
			}

			if matched && distance < 0.35 {
				fmt.Printf("\tFROM %d: %s (%s)\n\tSIMILARITY: %0.2f\n\tDISTANCE: %0.2f\n\tTIMESTAMP: %s\n\n", comparePos, compare[comparePos].Content, compareText, similarity, distance, compareTimestamp)
				target[targetLine.Position-1].Timestamp = compareTimestamp
				target[targetLine.Position-1].TimestampInferred = false
				lastMatchedTimestamp = compareTimestamp
				numMatched++
				break
			}
		}
	}
	fmt.Printf("\nCOMPLETED with %d matched of %d (%0.2f%%)\n", numMatched, len(transcript), float64(numMatched)/float64(len(transcript))*100)

	return target
}

func distanceWithinNPcnt(x int, y int, pcnt float64) bool {
	return distancePcnt(x, y) <= pcnt
}

func distancePcnt(x int, y int) float64 {
	distance := math.Abs(float64(x - y))
	return distance / float64(x)
}

func cleanString(raw string) string {
	raw = strings.Replace(raw, "'s", "", -1)
	raw = strings.Replace(raw, "’s", "", -1)
	raw = strings.Replace(raw, "-", " ", -1)
	raw = punctuation.ReplaceAllString(raw, "")
	return strings.TrimSpace(stopwords.CleanString(strings.ToLower(raw), "en", false))
}

func assemblyAiToDialog(aai []*assemblyai.TranscriptUtterance) []models.Dialog {
	dialog := make([]models.Dialog, len(aai))
	for k, v := range aai {
		dialog[k] = models.Dialog{
			Type:      models.DialogTypeChat,
			Content:   v.Text,
			Actor:     v.Speaker, // A, B or C
			Timestamp: time.Duration(v.Start) * time.Millisecond,
			Position:  int64(k) + 1,
		}
	}
	return dialog
}
