package assemblyai

import (
	"fmt"
	"io"
	"math"
)

func ToFlatFile(rawData *TranscriptionStatusResponse, outputWriter io.Writer) error {
	for _, v := range rawData.Utterances {
		if _, err := fmt.Fprintf(outputWriter, "#OFFSET: %d\n", int64(math.Round(float64(v.Start)/1000))); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(outputWriter, "Unknown %s: %s\n", v.Speaker, v.Text); err != nil {
			return err
		}
	}
	return nil
}
