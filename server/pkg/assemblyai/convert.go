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

// ToSrt not tested
func ToSrt(rawData *TranscriptionStatusResponse, outputWriter io.Writer) error {
	for _, v := range rawData.Utterances {
		if _, err := fmt.Fprintf(outputWriter, "%d\n", v.Pos); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(
			outputWriter,
			"%s --> %s\n",
			formatDurationAsSrtTimestamp(time.Duration(v.Start)*time.Millisecond),
			formatDurationAsSrtTimestamp(time.Duration(v.End)*time.Millisecond),
		); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(outputWriter, "%s: %s\n", v.Speaker, v.Text); err != nil {
			return err
		}
		if _, err := fmt.Fprint(outputWriter, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func formatDurationAsSrtTimestamp(dur time.Duration) string {
	return time.Unix(0, 0).UTC().Add(dur).Format("15:04:05,000")
}
