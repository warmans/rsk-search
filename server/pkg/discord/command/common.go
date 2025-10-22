package command

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

var punctuation = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
var spaces = regexp.MustCompile(`[\s]{2,}`)
var metaWhitespace = regexp.MustCompile(`[\n\r\t]+`)

var extractState = regexp.MustCompile(`\|\|(\{.*\})\|\|`)

func mustEncodeState[T any](s T) string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("||%s||", string(b))
}

func decodeState[T any](raw string) (*T, error) {
	state := new(T)
	err := json.Unmarshal([]byte(strings.Trim(raw, "|")), state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func extractStateFromBody[T any](msg *discordgo.Message) (*T, error) {
	foundState := extractState.FindString(msg.Content)
	if foundState == "" {
		return nil, fmt.Errorf("failed to find state in message body")
	}

	state, err := decodeState[T](foundState)
	if err != nil {
		return nil, fmt.Errorf("failed to parse state: %s", foundState)
	}

	return state, nil
}

func deleteStateFromContent(content string) string {
	return string(extractState.ReplaceAll([]byte(content), []byte{}))
}

func contentToFilename(rawContent string) string {
	rawContent = punctuation.ReplaceAllString(rawContent, "")
	rawContent = spaces.ReplaceAllString(rawContent, " ")
	rawContent = metaWhitespace.ReplaceAllString(rawContent, " ")
	rawContent = strings.ToLower(strings.TrimSpace(rawContent))
	split := strings.Split(rawContent, " ")
	if len(split) > 9 {
		split = split[:8]
	}
	return strings.Join(split, "-")
}
