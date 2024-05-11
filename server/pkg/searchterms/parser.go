package searchterms

import (
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/filter"
	"strings"
)

type Term struct {
	Field string
	Value string
	Op    filter.CompOp
}

func MustParse(s string) []Term {
	f, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return f
}

func Parse(s string) ([]Term, error) {
	if s == "" {
		return nil, nil
	}
	return newParser(newScanner(s)).Parse()
}

func newParser(s *scanner) *parser {
	return &parser{s: s}
}

type parser struct {
	s      *scanner
	peeked *token
}

func (p *parser) Parse() ([]Term, error) {
	terms, err := p.parseOuter()
	if err != nil {
		return nil, err
	}
	if _, err := p.requireNext(tagEOF); err != nil {
		return nil, err
	}
	return terms, nil
}

func (p *parser) parseOuter() ([]Term, error) {
	terms := []Term{}
	term, err := p.parseInner()
	if err != nil {
		return nil, err
	}
	for term != nil {
		terms = append(terms, *term)
		term, err = p.parseInner()
		if err != nil {
			return nil, err
		}
	}
	return terms, nil
}

func (p *parser) parseInner() (*Term, error) {
	tok, err := p.getNext()
	if err != nil {
		return nil, err
	}
	switch tok.tag {
	case tagEOF:
		return nil, nil
	case tagQuotedString:
		return &Term{
			Field: "content",
			Value: strings.Trim(tok.lexeme, `"`),
			Op:    filter.CompOpEq,
		}, nil
	case tagWord:
		words := []string{tok.lexeme}
		next, err := p.peekNext()
		if err != nil {
			return nil, err
		}
		for next.tag == tagWord {
			next, err = p.getNext()
			if err != nil {
				return nil, err
			}
			if word := strings.TrimSpace(next.lexeme); word != "" {
				words = append(words, word)
			}
			next, err = p.peekNext()
			if err != nil {
				return nil, err
			}
		}
		return &Term{
			Field: "content",
			Value: strings.Join(words, " "),
			Op:    filter.CompOpFuzzyLike,
		}, nil
	case tagMention:
		mentionText, err := p.requireNext(tagQuotedString, tagWord, tagEOF)
		if err != nil {
			return nil, err
		}
		return &Term{
			Field: "actor",
			Value: strings.ToLower(mentionText.lexeme),
			Op:    filter.CompOpEq,
		}, nil
	case tagPublication:
		mentionText, err := p.requireNext(tagQuotedString, tagWord, tagEOF)
		if err != nil {
			return nil, err
		}
		return &Term{
			Field: "publication",
			Value: strings.ToLower(mentionText.lexeme),
			Op:    filter.CompOpEq,
		}, nil
	default:
		return nil, errors.Errorf("unexpected token '%s'", tok)
	}
}

// peekNext gets the next token without advancing.
func (p *parser) peekNext() (token, error) {
	if p.peeked != nil {
		return *p.peeked, nil
	}
	t, err := p.getNext()
	if err != nil {
		return token{}, err
	}
	p.peeked = &t
	return t, nil
}

// getNext advances to the next token
func (p *parser) getNext() (token, error) {
	if p.peeked != nil {
		t := *p.peeked
		p.peeked = nil
		return t, nil
	}
	t, err := p.s.next()
	if err != nil {
		return token{}, err
	}
	return t, err
}

// requireNext advances to the next token and asserts it is one of the given tags.
func (p *parser) requireNext(oneOf ...tag) (token, error) {
	t, err := p.getNext()
	if err != nil {
		return t, err
	}
	for _, tag := range oneOf {
		if t.tag == tag {
			return t, nil
		}
	}
	return token{}, errors.Errorf("expected one of '%v', found '%s'", oneOf, t.tag)
}
