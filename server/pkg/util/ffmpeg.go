package util

import "strings"

func FfmpegSanitizeDrawtext(text string) string {
	text = strings.Replace(text, ":", `\:`, -1)
	text = strings.Replace(text, "'", `\'`, -1)
	return text
}
