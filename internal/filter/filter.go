package filter

type CompOp string

const (
	CompOpEq   CompOp = "="
	CompOpNeq  CompOp = "!="
	CompOpLike CompOp = "~="
	CompOpLt   CompOp = "<"
	CompOpLe   CompOp = "<="
	CompOpGt   CompOp = ">"
	CompOpGe   CompOp = ">="
)

func (op CompOp) Precedence() int {
	return 3
}

type BoolOp string

const (
	BoolOpAnd BoolOp = "and"
	BoolOpOr  BoolOp = "or"
)

func (op BoolOp) Precedence() int {
	switch op {
	case BoolOpAnd:
		return 2
	case BoolOpOr:
		return 1
	default:
		return 0
	}
}

type Visitor interface {
	VisitCompFilter(*CompFilter) (Visitor, error)
	VisitBoolFilter(*BoolFilter) (Visitor, error)
}

type Filter interface {
	Accept(Visitor) error
	Precedence() int
}

type CompFilter struct {
	Field string
	Op    CompOp
	Value Value
}

func (c *CompFilter) Accept(visitor Visitor) error {
	if _, err := visitor.VisitCompFilter(c); err != nil {
		return err
	}
	return nil
}

func (c *CompFilter) Precedence() int {
	return c.Op.Precedence()
}

type BoolFilter struct {
	LHS Filter
	Op  BoolOp
	RHS Filter
}

func (b *BoolFilter) Accept(visitor Visitor) error {
	v, err := visitor.VisitBoolFilter(b)
	if err != nil {
		return err
	}
	if v == nil {
		return nil
	}
	if b.RHS != nil {
		if err := b.RHS.Accept(v); err != nil {
			return err
		}
	}
	if b.LHS != nil {
		if err := b.LHS.Accept(v); err != nil {
			return err
		}
	}
	return nil
}

func (b *BoolFilter) Precedence() int {
	return b.Op.Precedence()
}

func And(lhs, rhs Filter, filters ...Filter) Filter {
	filter := &BoolFilter{lhs, BoolOpAnd, rhs}
	for _, f := range filters {
		filter = &BoolFilter{filter, BoolOpAnd, f}
	}
	return filter
}

func Or(lhs, rhs Filter, filters ...Filter) Filter {
	filter := &BoolFilter{lhs, BoolOpOr, rhs}
	for _, f := range filters {
		filter = &BoolFilter{filter, BoolOpOr, f}
	}
	return filter
}

func Eq(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpEq, Value: val}
}

func Neq(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpNeq, Value: val}
}

func Gt(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpGt, Value: val}
}

func Ge(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpGe, Value: val}
}

func Lt(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpLt, Value: val}
}

func Le(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpLe, Value: val}
}

func Like(field string, val Value) Filter {
	return &CompFilter{Field: field, Op: CompOpLike, Value: val}
}
