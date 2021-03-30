package tscript

import (
	"bufio"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/warmans/rsk-search/pkg/models"
	"strconv"
	"strings"
	"unicode"
)

const PosSpacing = 100

// Import imports plain text transcripts to JSON.
func Import(scanner *bufio.Scanner, startPos int64) ([]models.Dialog, []models.Synopsis, error) {

	output := make([]models.Dialog, 0)
	position := startPos
	lastOffset := int64(0)
	numOffsets := 0

	synopsies := make([]models.Synopsis, 0)
	var currentSynopsis *models.Synopsis

	for scanner.Scan() {
		position += PosSpacing

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// OFFSET lines related to the next line of text so just store the offset
		// and continue.
		if IsOffsetTag(line) {
			if offset, ok := ScanOffset(line); ok {
				if offset <= lastOffset {
					return nil, nil, fmt.Errorf("offsets are invalid")
				}
				lastOffset = offset
				numOffsets++
			}
			continue
		}

		if strings.HasPrefix(line, "#SYN: ") || strings.HasPrefix(line, "#/SYN") {
			if currentSynopsis != nil {
				currentSynopsis.EndPos = position-PosSpacing
				synopsies = append(synopsies, *currentSynopsis)
				currentSynopsis = nil
			}
			if strings.HasPrefix(line, "#SYN: ") {
				currentSynopsis = &models.Synopsis{Description: CorrectContent(strings.TrimSpace(strings.TrimPrefix(line, "#SYN:"))), StartPos: position}
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
			if actor != "none" {
				di.Actor = actor
			}
		}
		di.Content = CorrectContent(strings.TrimSpace(parts[1]))

		output = append(output, di)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	if numOffsets == 0 {
		return nil, nil, fmt.Errorf("document appears to be missing offsets")
	}

	if currentSynopsis != nil {
		currentSynopsis.EndPos = position
		synopsies = append(synopsies, *currentSynopsis)
	}

	return output, synopsies, nil
}

func CorrectContent(c string) string {
	runes := []rune(c)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
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
