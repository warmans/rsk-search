package transcript

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/models"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type tag string

func (t tag) Open() string {
	return fmt.Sprintf("#%s:", string(t))
}

func (t tag) Close() string {
	return fmt.Sprintf("#/%s", string(t))
}

const (
	OffsetTag   tag = "OFFSET"
	GapTag      tag = "GAP"
	SynopsisTag tag = "SYN"
	TriviaTag   tag = "TRIVIA"
)

var ErrEOF = fmt.Errorf("EOF")

const PosSpacing = 1

type exportOptions struct {
	stripMetadata bool
}

type ExportOption func(opts *exportOptions)

func WithStripMetadata() ExportOption {
	return func(opts *exportOptions) {
		opts.stripMetadata = true
	}
}

func NewTranscriptScanner(scanner *bufio.Scanner) *TranscriptScanner {
	return &TranscriptScanner{scanner: scanner}
}

type TranscriptScanner struct {
	scanner *bufio.Scanner
	peeked  *string
}

func (ts *TranscriptScanner) Next() (string, error) {
	if peeked := ts.peeked; peeked != nil {
		ts.peeked = nil
		return *peeked, nil
	}
	if !ts.scanner.Scan() {
		return "", ErrEOF
	}
	return ts.scanner.Text(), nil
}

func (ts *TranscriptScanner) PeekNext() (string, error) {
	if ts.peeked != nil {
		return *ts.peeked, nil
	}
	peeked, err := ts.Next()
	if err != nil {
		return "", err
	}
	ts.peeked = &peeked
	return peeked, nil
}

func (ts *TranscriptScanner) Err() error {
	return ts.scanner.Err()
}

// ReadAllPrefixed reads all lines from the scanner that start with the given prefix.
// this is useful for reading synopsis and trivia lines.
func (ts *TranscriptScanner) ReadAllPrefixed(prefix string) ([]string, error) {
	all := []string{}
	for {
		peeked, err := ts.PeekNext()
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}
			return nil, err
		}
		// stop if another tag is encountered
		if IsTag(peeked) {
			return all, nil
		}
		if strings.HasPrefix(peeked, prefix) {
			next, err := ts.Next()
			if err != nil {
				return nil, err
			}
			all = append(all, strings.TrimSpace(strings.TrimPrefix(next, prefix)))
			continue
		}
		break
	}
	return all, nil
}

func Validate(scanner *bufio.Scanner) error {
	transcript, err := Import(scanner, "", 0)
	if err != nil {
		return err
	}
	if len(transcript.Transcript) == 0 {
		return fmt.Errorf("no valid lines parsed from transcript")
	}
	return nil
}

// Import imports plain text transcripts to JSON.
func Import(scanner *bufio.Scanner, episodeID string, startPos int64) (*models.Transcript, error) {

	parser := NewTranscriptScanner(scanner)

	position := startPos
	var lastOffset time.Duration
	var numOffsets int

	transcript := &models.Transcript{
		Transcript: make([]models.Dialog, 0),
		Synopsis:   make([]models.Synopsis, 0),
		Trivia:     make([]models.Trivia, 0),
	}

	var currentSynopsis *models.Synopsis
	var currentTrivia *models.Trivia

	for {
		notable := false

		currentLine, err := parser.Next()
		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}
			return nil, err
		}

		// strip space and non-breakable-spaces
		line := strings.TrimSpace(strings.ReplaceAll(currentLine, "\u00a0", " "))
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
			if offset, ok := ScanSeconds(OffsetTag, line); ok {
				if offset > 0 && offset <= lastOffset {
					return nil, fmt.Errorf("offsets are invalid")
				}
				lastOffset = offset
				numOffsets++
			}
			continue
		}

		if IsGapTag(line) {
			duration, err := ScanDuration(GapTag, line)
			if err == nil {
				gap := models.Dialog{
					ID:       models.DialogID(episodeID, position),
					Type:     models.DialogTypeGap,
					Position: position,
					Duration: duration,
				}
				if lastOffset > 0 {
					gap.Timestamp = lastOffset
					gap.TimestampInferred = false
					lastOffset = 0
				}
				transcript.Transcript = append(transcript.Transcript, gap)
				position += PosSpacing
			}
			continue
		}

		if IsSynopsisTag(line) {
			if currentSynopsis != nil {
				currentSynopsis.EndPos = position
				transcript.Synopsis = append(transcript.Synopsis, *currentSynopsis)
				currentSynopsis = nil
			}
			if strings.HasPrefix(line, "#SYN: ") {
				currentSynopsis = &models.Synopsis{Description: CorrectContent(strings.TrimSpace(strings.TrimPrefix(line, "#SYN:"))), StartPos: position}
				nextLines, err := parser.ReadAllPrefixed("#")
				if err != nil {
					return nil, err
				}
				if len(nextLines) > 0 {
					currentSynopsis.Description += "\n" + strings.Join(nextLines, "\n")
				}
			}
			continue
		}
		if IsTriviaTag(line) {
			if currentTrivia != nil {
				currentTrivia.EndPos = position
				transcript.Trivia = append(transcript.Trivia, *currentTrivia)
				currentTrivia = nil
			}
			if strings.HasPrefix(line, "#TRIVIA:") {
				currentTrivia = &models.Trivia{Description: CorrectContent(strings.TrimSpace(strings.TrimPrefix(line, "#TRIVIA:"))), StartPos: position}
				nextLines, err := parser.ReadAllPrefixed("#")
				if err != nil {
					return nil, err
				}
				if len(nextLines) > 0 {
					currentTrivia.Description += "\n" + strings.Join(nextLines, "\n")
				}
			}
			continue
		}

		di := models.Dialog{
			ID:                models.DialogID(episodeID, position),
			Type:              models.DialogTypeUnknown,
			Position:          position,
			Notable:           notable,
			TimestampInferred: true,
		}
		if lastOffset > 0 {
			di.Timestamp = lastOffset
			di.TimestampInferred = false
			lastOffset = 0
		}

		// line should be in the format "actor: text..."
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Printf("WARN: invalid line detected (missing actor): %s\n", line)
			parts = []string{"none", line}
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

		transcript.Transcript = append(transcript.Transcript, di)

		// advance to next line
		position += PosSpacing
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if currentSynopsis != nil {
		currentSynopsis.EndPos = position
		transcript.Synopsis = append(transcript.Synopsis, *currentSynopsis)
	}
	if currentTrivia != nil {
		currentTrivia.EndPos = position
		transcript.Trivia = append(transcript.Trivia, *currentTrivia)
	}

	return transcript, nil
}

// Export dumps dialog back to the raw format.
// the problem is this loses most of the metadata and it's not at all easy to
func Export(dialog []models.Dialog, synopsis []models.Synopsis, trivia []models.Trivia, opts ...ExportOption) (string, error) {

	options := &exportOptions{}
	for _, v := range opts {
		v(options)
	}

	output := strings.Builder{}
	for _, d := range dialog {
		if !options.stripMetadata {

			if d.Type == models.DialogTypeGap {
				output.WriteString(fmt.Sprintf("#GAP: %s\n", d.Duration.String()))
				continue
			}

			if d.Timestamp > 0 && !d.TimestampInferred {
				output.WriteString(fmt.Sprintf("#OFFSET: %0.2f\n", d.Timestamp.Seconds()))
			}
			for _, syn := range synopsis {
				if d.Position == syn.StartPos {
					synopsisLines := strings.Split(syn.Description, "\n")
					output.WriteString(fmt.Sprintf("#SYN: %s\n", synopsisLines[0]))
					if len(synopsisLines) > 1 {
						for _, line := range synopsisLines[1:] {
							output.WriteString(fmt.Sprintf("# %s\n", strings.TrimSpace(line)))
						}
					}
				}
				if d.Position == syn.EndPos {
					output.WriteString("#/SYN\n")
				}
			}
			for _, triv := range trivia {
				if d.Position == triv.StartPos {
					triviaLines := strings.Split(triv.Description, "\n")
					output.WriteString(fmt.Sprintf("#TRIVIA: %s\n", triviaLines[0]))
					if len(triviaLines) > 1 {
						for _, line := range triviaLines[1:] {
							output.WriteString(fmt.Sprintf("# %s\n", strings.TrimSpace(line)))
						}
					}
				}
				if d.Position == triv.EndPos {
					output.WriteString("#/TRIVIA\n")
				}
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

func IsTag(line string) bool {
	return IsOffsetTag(line) || IsTriviaTag(line) || IsSynopsisTag(line)
}

func IsOffsetTag(line string) bool {
	return strings.HasPrefix(line, OffsetTag.Open())
}

func IsGapTag(line string) bool {
	return strings.HasPrefix(line, GapTag.Open())
}

func IsTriviaTag(line string) bool {
	return strings.HasPrefix(line, TriviaTag.Open()) || strings.HasPrefix(line, TriviaTag.Close())
}

func IsSynopsisTag(line string) bool {
	return strings.HasPrefix(line, SynopsisTag.Open()) || strings.HasPrefix(line, SynopsisTag.Close())
}

func ScanSeconds(tagPrefix tag, line string) (time.Duration, bool) {
	offsetStr := strings.TrimSpace(strings.TrimPrefix(line, tagPrefix.Open()))
	if off, err := strconv.ParseFloat(offsetStr, 64); err == nil {
		return time.Duration(off*1000) * time.Millisecond, true
	}
	return 0, false
}

func ScanDuration(tagPrefix tag, line string) (time.Duration, error) {
	strDuration := strings.TrimSpace(strings.TrimPrefix(line, tagPrefix.Open()))
	return time.ParseDuration(strDuration)
}
