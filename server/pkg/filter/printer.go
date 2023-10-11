package filter

import (
	"bytes"
	"fmt"
	"io"
)

func MustPrint(f Filter) string {
	if f == nil {
		return ""
	}
	s, err := Print(f)
	if err != nil {
		panic(err)
	}
	return s
}

func Print(filter Filter) (string, error) {
	buff := &bytes.Buffer{}
	if err := Fprint(filter, buff); err != nil {
		return "", err
	}
	return buff.String(), nil
}

func Fprint(filter Filter, writer io.Writer) error {
	v := &Printer{w: writer}
	return filter.Accept(v)
}

type Printer struct {
	w io.Writer
}

func (p *Printer) VisitCompFilter(filter *CompFilter) (Visitor, error) {

	if _, err := fmt.Fprintf(p.w, "%s %s %s", filter.Field, filter.Op, filter.Value.String()); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Printer) VisitBoolFilter(filter *BoolFilter) (Visitor, error) {

	needsLeftParen := filter.LHS.Precedence() < filter.Precedence()

	if needsLeftParen {
		if _, err := fmt.Fprint(p.w, `(`); err != nil {
			return nil, err
		}
	}
	if err := filter.LHS.Accept(p); err != nil {
		return nil, err
	}
	if needsLeftParen {
		if _, err := fmt.Fprint(p.w, `)`); err != nil {
			return nil, err
		}
	}

	if _, err := fmt.Fprintf(p.w, " %s ", string(filter.Op)); err != nil {
		return nil, err
	}

	needsRightParen := filter.RHS.Precedence() < filter.Precedence()

	if needsRightParen {
		if _, err := fmt.Fprint(p.w, `(`); err != nil {
			return nil, err
		}
	}
	if err := filter.RHS.Accept(p); err != nil {
		return nil, err
	}
	if needsRightParen {
		if _, err := fmt.Fprint(p.w, `)`); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
