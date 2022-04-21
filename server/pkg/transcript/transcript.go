package transcript

import (
	"bufio"
	"fmt"
	"github.com/warmans/rsk-search/pkg/models"
	"strconv"
	"strings"
	"unicode"
)

const PosSpacing = 1

func Validate(scanner *bufio.Scanner) error {
	lines, _, _, err := Import(scanner, "", 0)
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		return fmt.Errorf("no valid lines parsed from transcript")
	}
	return nil
}

// Import imports plain text transcripts to JSON.
func Import(scanner *bufio.Scanner, episodeID string, startPos int64) ([]models.Dialog, []models.Synopsis, []models.Trivia, error) {

	output := make([]models.Dialog, 0)
	position := startPos
	lastOffset := int64(0)
	numOffsets := 0

	synopsies := make([]models.Synopsis, 0)
	var currentSynopsis *models.Synopsis

	trivia := make([]models.Trivia, 0)
	var currentTrivia *models.Trivia

	for scanner.Scan() {
		notable := false

		// strip space and non-breakable-spaces
		line := strings.TrimSpace(strings.ReplaceAll(scanner.Text(), "\u00a0", " "))
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
					return nil, nil, nil, fmt.Errorf("offsets are invalid")
				}
				lastOffset = offset
				numOffsets++
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
				currentSynopsis = &models.Synopsis{Description: CorrectContent(strings.TrimSpace(strings.TrimPrefix(line, "#SYN:"))), StartPos: position}
			}
			continue
		}
		if strings.HasPrefix(line, "#TRIVIA: ") || strings.HasPrefix(line, "#/TRIVIA") {
			if currentTrivia != nil {
				currentTrivia.EndPos = position
				trivia = append(trivia, *currentTrivia)
				currentTrivia = nil
			}
			if strings.HasPrefix(line, "#TRIVIA: ") {
				currentTrivia = &models.Trivia{Description: CorrectContent(strings.TrimSpace(strings.TrimPrefix(line, "#TRIVIA:"))), StartPos: position}
			}
			continue
		}

		position += PosSpacing
		di := models.Dialog{
			ID:             models.DialogID(episodeID, position),
			Type:           models.DialogTypeUnkown,
			Position:       position,
			Notable:        notable,
			OffsetInferred: true,
		}
		if lastOffset > 0 {
			di.OffsetSec = lastOffset
			di.OffsetInferred = false
			lastOffset = 0
		}

		// line should be in the format "actor: text..."
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, nil, nil, fmt.Errorf("line did not start with actor name or tag: %s", line)
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
		return nil, nil, nil, err
	}
	if currentSynopsis != nil {
		currentSynopsis.EndPos = position
		synopsies = append(synopsies, *currentSynopsis)
	}
	if currentTrivia != nil {
		currentTrivia.EndPos = position
		trivia = append(trivia, *currentTrivia)
	}

	return output, synopsies, trivia, nil
}

// Export dumps dialog back to the raw format.
// the problem is this loses most of the metadata and it's not at all easy to
func Export(dialog []models.Dialog, synopsis []models.Synopsis, trivia []models.Trivia) (string, error) {

	output := strings.Builder{}
	for _, d := range dialog {
		if d.OffsetSec > 0 && !d.OffsetInferred {
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
		for _, triv := range trivia {
			if d.Position == triv.StartPos {
				output.WriteString(fmt.Sprintf("#TRIVIA: %s\n", triv.Description))
			}
			if d.Position == triv.EndPos {
				output.WriteString("#/TRIVIA\n")
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
