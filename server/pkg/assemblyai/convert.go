package assemblyai

import (
	"fmt"
	"io"
	"time"
)

func ToFlatFile(rawData *TranscriptionStatusResponse, outputWriter io.Writer) error {
	for _, v := range rawData.Utterances {
		if _, err := fmt.Fprintf(outputWriter, "#OFFSET: %0.2f\n", (time.Duration(v.Start) * time.Millisecond).Seconds()); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(outputWriter, "Unknown %s: %s\n", v.Speaker, v.Text); err != nil {
			return err
		}
	}
	return nil
}
