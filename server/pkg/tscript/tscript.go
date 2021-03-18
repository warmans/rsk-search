package tscript

import (
	"bufio"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"os"
	"strconv"
	"strings"
)

// Import imports plain text transcripts to JSON.
func Import(f *os.File) ([]models.Dialog, []models.Synopsis, error) {

	output := make([]models.Dialog, 0)
	position := int64(0)
	lastOffset := int64(0)

	synopsies := make([]models.Synopsis, 0)
	var currentSynopsis *models.Synopsis

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		position += 100

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// OFFSET lines related to the next line of text so just store the offset
		// and continue.
		if IsOffsetTag(line) {
			if offset, ok := ScanOffset(line); ok {
				lastOffset = offset
			}
			continue
		}

		if strings.HasPrefix(line, "#SYN: ") || strings.HasPrefix(line, "#/SYN") {
			if currentSynopsis != nil {
				currentSynopsis.EndPos = position
				synopsies = append(synopsies, *currentSynopsis)
				currentSynopsis = nil
			}
			if strings.HasPrefix(line, "#SYN: ") {
				currentSynopsis = &models.Synopsis{Description: strings.TrimSpace(strings.TrimPrefix(line, ":")), StartPos: position}
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
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("line did not start with actor name or tag: %s", line)
		}

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
		return nil, nil, err
	}

	if currentSynopsis != nil {
		currentSynopsis.EndPos = position
		synopsies = append(synopsies, *currentSynopsis)
	}

	return output, synopsies, nil
}

func IsOffsetTag(line string) bool {
	return strings.HasPrefix(line, "#OFFSET:")
}

func ScanOffset(line string) (int64, bool) {
	offsetStr := strings.TrimSpace(strings.TrimPrefix(line, "#OFFSET:"))
	if off, err := strconv.Atoi(offsetStr); err == nil {
		return int64(off), true
	}
	return 0, false
}
