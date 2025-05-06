package data

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/warmans/rsk-search/pkg/assemblyai"
	"github.com/warmans/rsk-search/pkg/data"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"math"
	"path"
	"time"
)

// MergeAAISentimentsCommand
// e.g.  make build && ./bin/rsk-search data merge-aai-sentiments -s "var/aai-transcripts/xfm-S1E03.mp3?remastered=1.json" -t "ep-xfm-S1E03.json"
func MergeAAISentimentsCommand() *cobra.Command {
	var aaiDataPath string
	var targetTranscriptName string

	cmd := &cobra.Command{
		Use:   "merge-aai-sentiments",
		Short: "merge aai sentiments into transcript",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetPath := path.Join(cfg.dataDir, targetTranscriptName)
			target, err := data.LoadEpisodePath(targetPath)
			if err != nil {
				return fmt.Errorf("failed to load target transcript at path %s: %w", targetPath, err)
			}

			aaiPayload := assemblyai.TranscriptionStatusResponse{}
			if err := util.WithReadJSONFileDecoder(aaiDataPath, func(dec *json.Decoder) error {
				return dec.Decode(&aaiPayload)

			}); err != nil {
				return fmt.Errorf("failed to read timestamp source %s: %w", aaiDataPath, err)
			}

			return data.ReplaceEpisodeFile(cfg.dataDir, mergeSentiments(target, aaiPayload))
		},
	}

	cmd.Flags().StringVarP(&aaiDataPath, "aai-response", "s", "", "Dump of assembly AI response for episode")
	cmd.Flags().StringVarP(&targetTranscriptName, "target-transcript", "t", "", "Target transcript")

	return cmd
}

func mergeSentiments(target *models.Transcript, aaiResponse assemblyai.TranscriptionStatusResponse) *models.Transcript {
	if aaiResponse.SentimentAnalysisResults == nil {
		return target
	}

	for k, dialog := range target.Transcript {
		sentiments := []*assemblyai.Sentiment{}
		for _, s := range aaiResponse.SentimentAnalysisResults {
			if s.Sentiment == "NEUTRAL" {
				continue
			}
			dialogEnd := dialog.Timestamp + dialog.Duration
			sentimentStart := time.Duration(s.Start) * time.Millisecond
			sentimentEnd := time.Duration(s.End) * time.Millisecond

			if sentimentStart >= dialog.Timestamp && sentimentEnd <= dialogEnd {
				// sentiment is contained within the dialog line
				sentiments = append(sentiments, s)
			} else {
				//sentiment overlaps somehow
				overlapStart := math.Max(float64(dialog.Timestamp), float64(sentimentStart))
				overlapEnd := math.Min(float64(dialogEnd), float64(sentimentEnd))
				overlapPercent := (overlapEnd - overlapStart) / float64(dialog.Duration)

				// only accept sentiments that overlap by more than 30% with the original dialog line
				if overlapPercent > 0.30 {
					sentiments = append(sentiments, s)
				}
			}
			target.Transcript[k].Sentiment = overallSentiment(sentiments)
		}
	}

	return target
}

func overallSentiment(sentiments []*assemblyai.Sentiment) models.Sentiment {
	if len(sentiments) == 0 {
		return models.SentimentMixed
	}
	var positive int
	var negative int
	for _, s := range sentiments {
		switch s.Sentiment {
		case "POSITIVE":
			positive++
		case "NEGATIVE":
			negative++
		}
	}
	if positive == negative {
		return models.SentimentMixed
	}
	if positive > negative {
		return models.SentimentPositive
	}

	return models.SentimentNegative
}
