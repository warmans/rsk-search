package searchterms

import (
	"fmt"
	"strings"
	"unicode"
)

type tag string

const (
	tagEOF tag = "EOF"

	tagMention     = "@"
	tagPublication = "~"

	tagQuotedString = "QUOTED_STRING"
	tagWord         = "WORD"
)

type token struct {
	tag    tag
	lexeme string
}

func (t token) String() string {
	return fmt.Sprintf("{%s: '%s'}", string(t.tag), t.lexeme)
}

func Scan(str string) ([]token, error) {

	scanner := newScanner(str)

	tokens := make([]token, 0)
	for {
		tok, err := scanner.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)

		if tok.tag == tagEOF {
			break
		}
	}
	return tokens, nil
}

func newScanner(str string) *scanner {
	return &scanner{
		input:  []rune(str),
		pos:    0,
		offset: 0,
	}
}

type scanner struct {
	input  []rune
	pos    int
	offset int
}

// Next gets the next token, advancing the scanner.
func (s *scanner) Next() (token, error) {
	return s.next()
}

func (s *scanner) next() (token, error) {
	s.skipWhitespace()
	if s.atEOF() {
		return s.emit(tagEOF), nil
	}
	switch r := s.nextRune(); r {
	case '@':
		return s.emit(tagMention), nil
	case '~':
		return s.emit(tagPublication), nil
	case '"':
		return s.scanString()
	default:
		if isValidInputRune(r) {
			return s.scanWord()
		}
		return s.error("unknown entity")
	}
}

func (s *scanner) nextRune() rune {
	r := s.input[s.pos]
	s.pos++
	return r
}

// matchNextRune will match the next rune of a multi-run tag e.g. >= !=
func (s *scanner) matchNextRune(r rune) bool {
	if s.atEOF() || s.peekRune() != r {
		return false
	}
	s.nextRune()
	return true
}

func (s *scanner) skipWhitespace() {
	for !s.atEOF() && unicode.IsSpace(s.peekRune()) {
		s.nextRune()
	}
	s.offset = s.pos
}

func (s *scanner) scanString() (token, error) {
	for !s.matchNextRune('"') {
		if s.atEOF() {
			return s.error("unclosed double quote")
		}
		s.nextRune()
	}
	return trimTokenLexeme(s.emit(tagQuotedString), `""`), nil
}

func (s *scanner) scanWord() (token, error) {
	for !s.atEOF() && isValidInputRune(s.peekRune()) && !isWhitespace(s.peekRune()) {
		s.nextRune()
	}
	return s.emit(tagWord), nil
}

func (s *scanner) atEOF() bool {
	return s.pos >= len(s.input)
}

func (s *scanner) peekRune() rune {
	return s.input[s.pos]
}

func (s *scanner) emit(tag tag) token {
	lexeme := string(s.input[s.offset:s.pos])
	s.offset = s.pos
	return token{tag: tag, lexeme: lexeme}
}

func (s *scanner) error(reason string) (token, error) {
	return token{}, fmt.Errorf("failed to scan string at position %d ('%s'): %s", s.pos, string(s.input[s.offset:s.pos]), reason)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isValidInputRune(r rune) bool {
	return r != '@' && r != '~' && r != '"'
}

func trimTokenLexeme(t token, trimSet string) token {
	t.lexeme = strings.Trim(t.lexeme, trimSet)
	return t
}
