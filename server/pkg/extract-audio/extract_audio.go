package extract_audio

import (
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"time"
)

func ExtractAudio(output io.Writer, inputFilePath string, fromTimestamp, toTimestamp time.Duration) error {
	return ffmpeg_go.
		Input(inputFilePath).
		Output("pipe:",
			ffmpeg_go.KwArgs{
				"format": "mp3",
				"vcodec": "copy",
				"acodec": "copy",
				"ss":     fmt.Sprintf("%0.2f", fromTimestamp.Seconds()),
				"to":     fmt.Sprintf("%0.2f", toTimestamp.Seconds()),
			},
		).WithOutput(output, os.Stderr).Run()
}
