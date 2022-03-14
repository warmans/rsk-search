package filter

import (
	"github.com/pkg/errors"
	"strconv"
)

func MustParse(s string) Filter {
	f, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return f
}

func Parse(s string) (Filter, error) {
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

func (p *parser) Parse() (Filter, error) {
	filter, err := p.parseOuter(1)
	if err != nil {
		return nil, err
	}
	if _, err := p.requireNext(tagEOF); err != nil {
		return nil, err
	}
	return filter, nil
}

func (p *parser) parseOuter(minPrec int) (Filter, error) {

	filter, err := p.parseInner()
	if err != nil {
		return nil, err
	}
	for {
		nextToken, err := p.peekNext()
		if err != nil {
			return nil, err
		}
		if !isBoolOp(nextToken) || precedence(nextToken.tag) < minPrec {
			break
		}

		op, err := p.getNext()
		if err != nil {
			return nil, err
		}
		rhs, err := p.parseOuter(precedence(op.tag) + 1)
		if err != nil {
			return nil, err
		}

		switch op.tag {
		case tagAnd:
			filter = And(filter, rhs)
		case tagOr:
			filter = Or(filter, rhs)
		default:
			panic(errors.Errorf("unexpected token '%s'", op.tag))
		}
	}

	return filter, nil
}

func (p *parser) parseInner() (Filter, error) {
	token, err := p.getNext()
	if err != nil {
		return nil, err
	}
	switch token.tag {
	case tagLParen:
		filter, err := p.parseOuter(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.requireNext(tagRParen); err != nil {
			return nil, err
		}
		return filter, nil
	case tagField:
		op, err := p.requireNext(tagEq, tagNeq, tagLt, tagLe, tagGe, tagGt, tagLike, tagFuzzy)
		if err != nil {
			return nil, err
		}
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		switch op.tag {
		case tagEq:
			return Eq(token.lexeme, val), nil
		case tagNeq:
			return Neq(token.lexeme, val), nil
		case tagLt:
			return Lt(token.lexeme, val), nil
		case tagLe:
			return Le(token.lexeme, val), nil
		case tagGt:
			return Gt(token.lexeme, val), nil
		case tagGe:
			return Ge(token.lexeme, val), nil
		case tagLike:
			return Like(token.lexeme, val), nil
		case tagFuzzy:
			return FuzzyLike(token.lexeme, val), nil
		default:
			panic(errors.Errorf("unexpected field token '%s'", op.tag))
		}
	default:
		return nil, errors.Errorf("unexpected token '%s'", token)
	}
}

func (p *parser) parseValue() (Value, error) {
	token, err := p.s.Next()
	if err != nil {
		return nil, err
	}
	switch token.tag {
	case tagNull:
		return Null(), nil
	case tagInt:
		i, err := strconv.ParseInt(token.lexeme, 10, 64)
		if err != nil {
			panic(errors.Wrapf(err, "failed to parse %s", token.tag))
		}
		return Int(i), nil
	case tagBool:
		b, err := strconv.ParseBool(token.lexeme)
		if err != nil {
			panic(errors.Wrapf(err, "failed to parse %s'", token))
		}
		if b {
			return Bool(true), nil
		}
		return Bool(false), nil
	case tagString:
		return String(token.lexeme), nil
	case tagFloat:
		f, err := strconv.ParseFloat(token.lexeme, 64)
		if err != nil {
			panic(errors.Wrapf(err, "failed to parse %s", token))
		}
		return Float(f), nil

	default:
		return nil, errors.Errorf("unexpected value tag %s", token)
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

func isBoolOp(token token) bool {
	return token.tag == tagAnd || token.tag == tagOr
}

func precedence(tag tag) int {
	switch tag {
	case tagAnd:
		return 2
	case tagOr:
		return 1
	default:
		return 0
	}
}
