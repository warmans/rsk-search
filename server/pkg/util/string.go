package util

import "strings"

func StripNonAlphanumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func LastSegment(str string, sep string) string {
	split := strings.Split(str, sep)
	return split[len(split)-1]
}
