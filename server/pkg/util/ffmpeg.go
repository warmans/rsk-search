package util

import "strings"

func FfmpegSanitizeDrawtext(text string) string {
	text = strings.ReplaceAll(text, ":", `\:`)
	text = strings.ReplaceAll(text, "'", `\'`)
	return text
}
