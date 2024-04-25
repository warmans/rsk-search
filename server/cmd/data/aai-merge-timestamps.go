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

func MergeTimestampsAAICommand() *cobra.Command {

	var timestampSourceFile string
	var targetTranscriptName string
	var outputPath string
	var replace bool
	var debugPos int64
	var debugComparePos int64

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

			target.Transcript = mergeTimestampsTo(target.Transcript, assemblyAiToDialog(timestampSource.Utterances), debugPos, debugComparePos)
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
	cmd.Flags().BoolVarP(&replace, "replace", "r", false, "replace source file")
	cmd.Flags().Int64VarP(&debugPos, "debug-pos", "p", 0, "Dump debug info for this position in the target transcript")
	cmd.Flags().Int64VarP(&debugComparePos, "debug-compare-pos", "c", 0, "Limit debug output to comparison lines with this position in the comparison transcript")

	return cmd
}

func mergeTimestampsTo(target []models.Dialog, compare []models.Dialog, debugPos int64, debugComparePos int64) []models.Dialog {

	// clear all non-chat data
	transcript := []models.Dialog{}
	for k, v := range target {
		//reset all timestamps
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
			standardMaxDistance := 0.040

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
	fmt.Printf("\nCOMPLETED with %d matched of %d (%0.2f%%)\n", numMatched, len(transcript), float64(numMatched)/float64(len(transcript))*100)

	return target
}

func distanceWithinNPcnt(x int, y int, total int, pcnt float64) bool {
	return distancePcnt(x, y, total) <= pcnt
}

func distancePcnt(x int, y int, total int) float64 {
	distance := math.Abs(float64(x - y))
	return distance / float64(total)
}

func cleanString(raw string) string {
	raw = strings.Replace(raw, "'s", "", -1)
	raw = strings.Replace(raw, "’s", "", -1)
	raw = strings.Replace(raw, "-", " ", -1)
	raw = strings.Replace(raw, "carl", "karl", -1)
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
