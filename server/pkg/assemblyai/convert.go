package assemblyai

import (
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"io"
)

// ToDialog will convert the assembly AI response to a standard transcript dialog.
func ToDialog(episodeID string, rawData *TranscriptionStatusResponse) ([]models.Dialog, error) {
	dialog := []models.Dialog{}
	for k, v := range rawData.Utterances {
		dialog = append(dialog, models.Dialog{
			ID:             models.DialogID(episodeID, int64(k)),
			Position:       int64(k),
			OffsetSec:      v.Start / 1000,
			OffsetInferred: false,
			Type:           models.DialogTypeChat,
			Actor:          v.Speaker,
			Content:        v.Text,
		})
	}
	return dialog, nil
}

func ToFlatFile(rawData *TranscriptionStatusResponse, outputWriter io.Writer) error {
	for _, v := range rawData.Utterances {

		if _, err := fmt.Fprintf(outputWriter, "#OFFSET: %d\n", v.Start/1000); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(outputWriter, "Unknown %s: %s\n", v.Speaker, v.Text); err != nil {
			return err
		}
	}
	return nil
}
