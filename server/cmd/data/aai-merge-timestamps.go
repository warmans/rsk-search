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
	"slices"
	"strings"
	"time"
)

var punctuation = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)

// MergeAAITimestampsCommand
// e.g. export EP=S2E08; ./bin/rsk-search data transcribe-assembly-ai -i "https://scrimpton.com/dl/media/xfm-${EP}.mp3?remastered=1" && ./bin/rsk-search data aai-merge-timestamps -s "var/aai-transcripts/xfm-${EP}.mp3?remastered=1.json" -t ep-xfm-${EP}.json
func MergeAAITimestampsCommand() *cobra.Command {

	var timestampSourceFile string
	var targetTranscriptName string
	var outputPath string
	var replace bool
	var preserveTimestamps bool
	var debugPos int64
	var debugComparePos int64
	var skipPositions []int
	var forceReplace bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "aai-merge-timestamps",
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

			newTranscript, err := mergeTimestampsTo(
				target.Transcript,
				assemblyAiToDialog(timestampSource.Utterances),
				debugPos,
				debugComparePos,
				skipPositions,
				preserveTimestamps,
				forceReplace,
			)
			if err != nil {
				return fmt.Errorf("failed to update transcript: %w", err)
			}
			if dryRun {
				return nil
			}

			target.Transcript = newTranscript
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
	cmd.Flags().IntSliceVarP(&skipPositions, "skip-positions", "x", []int{}, "Skips the given offsets e.g. 1,2,3")
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "replace source file")
	cmd.Flags().Int64VarP(&debugPos, "debug-pos", "p", 0, "Dump debug info for this position in the target transcript")
	cmd.Flags().Int64VarP(&debugComparePos, "debug-compare-pos", "c", 0, "Limit debug output to comparison lines with this position in the comparison transcript")
	cmd.Flags().BoolVarP(&preserveTimestamps, "preserve-timestamps", "", true, "keep existing timestamps")
	cmd.Flags().BoolVarP(&forceReplace, "force-replace", "f", false, "replace even if new transcript has fewer timestamps")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Print result only, don't update anything")

	return cmd
}

func mergeTimestampsTo(target []models.Dialog, compare []models.Dialog, debugPos int64, debugComparePos int64, skip []int, preserveTimestamps bool, forceReplace bool) ([]models.Dialog, error) {

	initialNumOffsets := 0
	for _, v := range target {
		if v.Timestamp > 0 && !v.TimestampInferred {
			initialNumOffsets++
		}
	}

	// clear all non-chat data
	transcript := []models.Dialog{}
	preservedTimestamps := map[int]time.Duration{}
	for k, v := range target {
		//reset all timestamps
		if !preserveTimestamps {
			preservedTimestamps[k] = target[k].Timestamp
		}
		target[k].Timestamp = 0
		target[k].TimestampInferred = true
		if v.Type == models.DialogTypeChat && v.Actor != "" {
			transcript = append(transcript, v)
		}
	}

	lastMatchedTimestamp := time.Duration(0)
	numMatched := 0
	// each time something is matched adjust the distance slightly otherwise the distance inaccuracies compound
	// throughout the length of the transcript
	distanceModifier := float64(0)
	for targetPos, targetLine := range transcript {
		targetText := cleanString(targetLine.Content)
		numTargetWords := len(strings.Split(targetText, " "))
		if numTargetWords <= 1 {
			if debugPos > 0 && debugPos == int64(targetPos) {
				fmt.Printf("\tSKIP (too few target words %d)\n", numTargetWords)
			}
			continue
		}
		if slices.Index(skip, targetPos) > -1 {
			fmt.Printf("\tSKIP (manually skipped)\n")
			continue
		}

		// consider explicit timestamps as matched lines
		if targetLine.Timestamp > 0 && !targetLine.TimestampInferred {
			lastMatchedTimestamp = targetLine.Timestamp
		}
		fmt.Printf("TARGET %d: %s (%s)\n", targetPos, targetLine.Content, targetText)
		for comparePos := 0; comparePos < len(compare); comparePos++ {
			compareText := cleanString(compare[comparePos].Content)
			compareTimestamp := compare[comparePos].Timestamp
			numCompareWords := len(strings.Split(compareText, " "))
			debug := (debugPos > 0 && debugPos == int64(targetPos)) && (debugComparePos == 0 || int64(comparePos) == debugComparePos-1)

			// do not backtrack
			if compareTimestamp < lastMatchedTimestamp {
				if debug {
					fmt.Printf("\tSKIP (compare ts %s less than last matched %s)\n", compareTimestamp, lastMatchedTimestamp)
				}
				continue
			}

			// skip obviously useless data
			if numCompareWords <= 1 {
				if debug {
					fmt.Printf("\tSKIP (too few compare words %d): %s (%s)\n", numCompareWords, compare[comparePos].Content, compareText)
				}
				continue
			}

			distance := math.Abs(distancePcnt(targetPos, comparePos, len(target)) - distanceModifier)
			var matched bool
			var similarity float64
			standardMaxDistance := 0.020

			if numTargetWords <= 3 {
				// compare the original text instead
				similarity = calculateSimilarity(targetLine.Content, compare[comparePos].Content)
				matched = similarity >= 0.60 && distance < 0.010 // use lower distance for short matches
			} else {
				compareText = makeComparisonSameLengthAsTarget(targetText, compareText, peekNextWords(compare, comparePos, numTargetWords)...)
				if numTargetWords > numCompareWords {
					// extend comparison with following lines to match target length
					similarity = calculateSimilarity(targetText, compareText)
					matched = similarity >= 0.65 && distance < standardMaxDistance
				} else {
					//truncate comparison to same length as target
					similarity = calculateSimilarity(targetText, compareText)
					matched = similarity >= 0.60 && distance < standardMaxDistance
				}
			}

			if matched || debug {
				fmt.Printf("\tFROM %d: %s (%s)\n\tSIMILARITY: %0.2f\n\tDISTANCE: %0.3f/%0.3f\n\tTIMESTAMP: %s\n\n", comparePos, compare[comparePos].Content, compareText, similarity, distance, distanceModifier, compareTimestamp)
			}
			if matched {
				target[targetLine.Position-1].Timestamp = compareTimestamp
				target[targetLine.Position-1].TimestampInferred = false
				lastMatchedTimestamp = compareTimestamp
				numMatched++
				distanceModifier = distancePcnt(targetPos, comparePos, len(target))
				break
			}
		}
	}

	// re-add old timestamps
	for pos, ts := range preservedTimestamps {
		if target[pos].TimestampInferred {
			target[pos].Timestamp = ts
			target[pos].TimestampInferred = false
		}
	}

	originalAccuracy := float64(initialNumOffsets) / float64(len(transcript)) * 100
	newAccuracy := float64(numMatched) / float64(len(transcript)) * 100

	if originalAccuracy > newAccuracy {
		fmt.Printf("\nWARNING: original accuracy %0.2f%% > new accuracy %0.2f%%\n", originalAccuracy, newAccuracy)
		if !forceReplace {
			return nil, fmt.Errorf("new accuracy %0.2f%% is lower than original accuracy %0.2f%%", newAccuracy, originalAccuracy)
		}
	}

	fmt.Printf("\nCOMPLETED with %d matched of %d (%0.2f -> %0.2f%%)\n", numMatched, len(transcript), originalAccuracy, newAccuracy)

	return target, nil
}

func distanceWithinNPcnt(x int, y int, total int, pcnt float64) bool {
	return distancePcnt(x, y, total) <= pcnt
}

func distancePcnt(x int, y int, total int) float64 {
	distance := math.Abs(float64(x - y))
	return distance / float64(total)
}

func cleanString(raw string) string {
	raw = strings.ReplaceAll(raw, "'s", "")
	raw = strings.ReplaceAll(raw, "’s", "")
	raw = strings.ReplaceAll(raw, "-", " ")
	raw = strings.ReplaceAll(raw, "carl", "karl")
	raw = punctuation.ReplaceAllString(raw, "")
	raw = stopwords.CleanString(strings.ToLower(raw), "en", false)
	raw = withoutWords(raw, "yeah", "sure", "hello", "alright", "dont", "know", "i", "uh", "huh")
	return strings.TrimSpace(raw)
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

func makeComparisonSameLengthAsTarget(targetString string, compareString string, paddingWords ...string) string {
	wantNumWords := len(strings.Split(targetString, " "))

	compare := strings.Split(compareString, " ")
	if len(compare) == wantNumWords {
		return compareString
	}
	if len(compare) > wantNumWords {
		return strings.Join(compare[:wantNumWords], " ")
	}
	if len(compare) < wantNumWords {
		for _, paddingWord := range paddingWords {
			compare = append(compare, paddingWord)
			if len(compare) == wantNumWords {
				return strings.Join(compare, " ")
			}
		}
	}
	// didn't make it to intended length
	return strings.Join(compare, " ")
}

func peekNextWords(dialog []models.Dialog, start, targetNumWords int) []string {
	if len(dialog) < start {
		return []string{}
	}
	peeked := []string{}
	for _, line := range dialog[start:] {
		for _, v := range strings.Split(line.Content, " ") {
			peeked = append(peeked, cleanString(v))
			if len(peeked) == targetNumWords {
				return peeked
			}
		}
	}
	return peeked
}

func withoutWords(str string, words ...string) string {
	out := []string{}
	for _, v := range strings.Split(str, " ") {
		if slices.Index(words, v) == -1 {
			out = append(out, v)
		}
	}
	return strings.Join(out, " ")
}

func calculateSimilarity(a, b string) float64 {
	return strutil.Similarity(a, b, metrics.NewJaccard())
}
