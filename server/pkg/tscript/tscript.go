package tscript

import (
	"bufio"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"os"
	"strconv"
	"strings"
)

// Import imports plain text transcripts to JSON.
func Import(f *os.File) ([]models.Dialog, error) {

	output := make([]models.Dialog, 0)
	position := int64(0)
	lastOffset := int64(0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		position += 100

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// OFFSET lines related to the next line of text so just store the offset
		// and continue.
		if strings.HasPrefix(line, "#OFFSET:") {
			offset, ok := scanOffset(line)
			if !ok {
				lastOffset = offset
			}
			continue
		}

		di := models.Dialog{
			ID:       shortuuid.New(),
			Type:     models.DialogTypeUnkown,
			Position: position,
		}
		if lastOffset > 0 {
			di.OffsetSec = lastOffset
			lastOffset = 0
		}

		// line should be in the format "actor: text..."
		parts := strings.SplitN(line, ":", 1)

		actor := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(parts[0]), ":"))
		if actor == "song" {
			di.Type = models.DialogTypeSong
		} else {
			di.Type = models.DialogTypeChat
			di.Actor = actor
		}
		di.Content = strings.TrimSpace(parts[1])

		output = append(output, di)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func scanOffset(line string) (int64, bool) {
	offsetStr := strings.TrimSpace(strings.TrimPrefix(line, "#OFFSET:"))
	if off, err := strconv.Atoi(offsetStr); err == nil {
		return int64(off), true
	}
	return 0, false
}
