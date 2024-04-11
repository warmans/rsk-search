package extract_audio

import (
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

func ExtractAudio(output io.Writer, inputFilePath string, fromTimestamp, toTimestamp int64) error {
	return ffmpeg_go.
		Input(inputFilePath).
		Output("pipe:",
			ffmpeg_go.KwArgs{
				"format": "mp3",
				"vcodec": "copy",
				"acodec": "copy",
				"ss":     fmt.Sprint(fromTimestamp),
				"to":     fmt.Sprint(toTimestamp),
			},
		).WithOutput(output, os.Stderr).Run()
}
