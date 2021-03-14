package bleve_query

import (
	"fmt"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/warmans/rsk-search/pkg/filter"
	"strings"
)

func FilterToQuery(f filter.Filter) (query.Query, error) {
	if f == nil {
		return query.NewMatchAllQuery(), nil
	}
	q := NewBleveQuery()
	if err := f.Accept(q); err != nil {
		return nil, err
	}
	return q.q, nil
}

func NewBleveQuery() *BleveQuery {
	return &BleveQuery{q: query.NewBooleanQuery([]query.Query{}, []query.Query{}, []query.Query{})}
}

type BleveQuery struct {
	q *query.BooleanQuery
}

func (j *BleveQuery) VisitCompFilter(f *filter.CompFilter) (filter.Visitor, error) {

	cond, err := j.condition(f.Field, f.Op, f.Value)
	if err != nil {
		return nil, err
	}
	j.q.AddMust(cond)

	return nil, nil
}

func (j *BleveQuery) VisitBoolFilter(f *filter.BoolFilter) (filter.Visitor, error) {

	q := query.NewBooleanQuery([]query.Query{}, []query.Query{}, []query.Query{})

	if err := f.LHS.Accept(j); err != nil {
		return nil, err
	}

	lhs := j.q
	j.q = query.NewBooleanQuery([]query.Query{}, []query.Query{}, []query.Query{})

	if err := f.RHS.Accept(j); err != nil {
		return nil, err
	}
	if f.Op == filter.BoolOpAnd {
		q.AddMust(lhs, j.q)
	}
	if f.Op == filter.BoolOpOr {
		return nil, fmt.Errorf("OR not supported by index")
	}

	j.q = q
	return nil, nil
}

func (j *BleveQuery) condition(field string, op filter.CompOp, value filter.Value) (query.Query, error) {

	switch op {
	case filter.CompOpEq:
		if value.Type() == filter.IntType {
			q := query.NewNumericRangeInclusiveQuery(floatP(float64(value.Value().(int64))), floatP(float64(value.Value().(int64))), boolP(true), boolP(true))
			q.SetField(field)
			return q, nil
		}
		q := query.NewMatchPhraseQuery(stripQuotes(value.String()))
		q.SetField(field)
		return q, nil
	case filter.CompOpLike:
		q := query.NewMatchQuery(stripQuotes(value.String()))
		q.SetField(field)
		q.SetFuzziness(0)
		return q, nil
	case filter.CompOpNeq:
		q := query.NewMatchPhraseQuery(stripQuotes(value.String()))
		q.SetField(field)
		return query.NewBooleanQueryForQueryString([]query.Query{}, []query.Query{}, []query.Query{q}), nil
	case filter.CompOpGt:
		switch value.Type() {
		case filter.IntType:
			q := query.NewNumericRangeQuery(floatP(float64(value.Value().(int64))), nil)
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := query.NewNumericRangeQuery(floatP(value.Value().(float64)), nil)
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the bleve mapping
			q := query.NewTermRangeQuery(stripQuotes(value.String()), "")
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	case filter.CompOpLt:
		switch value.Type() {
		case filter.IntType:
			q := query.NewNumericRangeQuery(nil, floatP(float64(value.Value().(int64))))
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := query.NewNumericRangeQuery(nil, floatP(value.Value().(float64)))
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the bleve mapping
			q := query.NewTermRangeQuery("", stripQuotes(value.String()))
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	case filter.CompOpGe:
		switch value.Type() {
		case filter.IntType:
			q := query.NewNumericRangeInclusiveQuery(floatP(float64(value.Value().(int64))), nil, boolP(true), nil)
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := query.NewNumericRangeInclusiveQuery(floatP(value.Value().(float64)), nil, boolP(true), nil)
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the bleve mapping
			q := query.NewTermRangeInclusiveQuery(stripQuotes(value.String()), "", boolP(true), nil)
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	case filter.CompOpLe:
		switch value.Type() {
		case filter.IntType:
			q := query.NewNumericRangeInclusiveQuery(nil, floatP(float64(value.Value().(int64))), nil, boolP(true))
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := query.NewNumericRangeInclusiveQuery(nil, floatP(value.Value().(float64)), nil, boolP(true))
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the bleve mapping
			q := query.NewTermRangeInclusiveQuery("", stripQuotes(value.String()), nil, boolP(true))
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	default:
		return nil, fmt.Errorf("operation %s was not implemented", string(op))
	}
}

func floatP(f float64) *float64 {
	return &f
}

func boolP(b bool) *bool {
	return &b
}

func stripQuotes(v string) string {
	return strings.Trim(v, `"`)
}
