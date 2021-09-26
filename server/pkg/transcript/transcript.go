package transcript

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
		notable := false

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// if the line starts with an exclamation, consider it a noteworthy quote.
		if strings.HasPrefix(line, "!") {
			notable = true
			line = strings.TrimPrefix(line, "!")
		}

		// OFFSET lines related to the next line of text so just store the offset
		// and continue.
		if IsOffsetTag(line) {
			if offset, ok := ScanOffset(line); ok {
				if offset > 0 && offset <= lastOffset {
					return nil, nil, fmt.Errorf("offsets are invalid")
				}
				lastOffset = offset
				numOffsets++
			}
			continue
		}

		if strings.HasPrefix(line, "#SYN: ") || strings.HasPrefix(line, "#/SYN") {
			if currentSynopsis != nil {
				currentSynopsis.EndPos = position - PosSpacing
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
			Notable:  notable,
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
	if currentSynopsis != nil {
		currentSynopsis.EndPos = position
		synopsies = append(synopsies, *currentSynopsis)
	}

	return output, synopsies, nil
}

// Export dumps dialog back to the raw format.
// the problem is this loses most of the metadata and it's not at all easy to
func Export(dialog []models.Dialog, synopsis []models.Synopsis) (string, error) {

	output := strings.Builder{}
	for _, d := range dialog {
		if d.OffsetSec > 0 && d.OffsetInferred == false {
			output.WriteString(fmt.Sprintf("#OFFSET: %d\n", d.OffsetSec))
		}
		for _, syn := range synopsis {
			if d.Position == syn.StartPos {
				output.WriteString(fmt.Sprintf("#SYN: %s\n", syn.Description))
			}
			if d.Position == syn.EndPos {
				output.WriteString("#/SYN\n")
			}
		}

		noteable := ""
		if d.Notable {
			noteable = "!"
		}
		actor := "none"
		switch d.Type {
		case models.DialogTypeChat:
			// if the actor isn't set just use none as it seems most "none" sections are marked as chat.
			if strings.TrimSpace(d.Actor) != "" {
				actor = d.Actor
			}
		case models.DialogTypeSong:
			actor = "song"
		}
		output.WriteString(fmt.Sprintf("%s%s: %s\n", noteable, actor, d.Content))
	}
	return output.String(), nil
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
