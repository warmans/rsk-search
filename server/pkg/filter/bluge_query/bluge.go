package bluge_query

import (
	"fmt"
	"github.com/blugelabs/bluge"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/search/v2/mapping"
	"math"
	"strings"
	"time"
)

func FilterToQuery(f filter.Filter) (bluge.Query, error) {
	if f == nil {
		return bluge.NewMatchAllQuery(), nil
	}
	q := NewBlugeQuery()
	if err := f.Accept(q); err != nil {
		return nil, err
	}
	return q.q, nil
}

func NewBlugeQuery() *BlugeQuery {
	return &BlugeQuery{q: bluge.NewBooleanQuery()}
}

type BlugeQuery struct {
	q *bluge.BooleanQuery
}

func (j *BlugeQuery) VisitCompFilter(f *filter.CompFilter) (filter.Visitor, error) {

	cond, err := j.condition(f.Field, f.Op, f.Value)
	if err != nil {
		return nil, err
	}
	j.q.AddMust(cond)

	return nil, nil
}

func (j *BlugeQuery) VisitBoolFilter(f *filter.BoolFilter) (filter.Visitor, error) {

	q := bluge.NewBooleanQuery()

	if err := f.LHS.Accept(j); err != nil {
		return nil, err
	}

	lhs := j.q
	j.q = bluge.NewBooleanQuery()

	if err := f.RHS.Accept(j); err != nil {
		return nil, err
	}
	if f.Op == filter.BoolOpAnd {
		q.AddMust(lhs, j.q)
	}
	if f.Op == filter.BoolOpOr {
		q.AddShould(lhs, j.q)
	}

	j.q = q
	return nil, nil
}

func (j *BlugeQuery) condition(field string, op filter.CompOp, value filter.Value) (bluge.Query, error) {

	switch op {
	case filter.CompOpEq:
		return j.eqFilter(field, value)
	case filter.CompOpNeq:
		q, err := j.eqFilter(field, value)
		if err != nil {
			return nil, err
		}
		return bluge.NewBooleanQuery().AddMustNot(q), nil
	case filter.CompOpLike:
		q := bluge.NewMatchQuery(stripQuotes(value.String()))
		q.SetField(field)
		q.SetFuzziness(0)
		return q, nil
	case filter.CompOpFuzzyLike:
		q := bluge.NewMatchQuery(stripQuotes(value.String()))
		q.SetField(field)
		q.SetFuzziness(1)
		return q, nil
	case filter.CompOpGt:
		switch value.Type() {
		case filter.IntType:
			// is max always required?
			q := bluge.NewNumericRangeQuery(float64(value.Value().(int64)), math.MaxFloat64)
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := bluge.NewNumericRangeQuery(value.Value().(float64), math.MaxFloat64)
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the document mapping
			q := bluge.NewTermRangeQuery(stripQuotes(value.String()), "")
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	case filter.CompOpLt:
		switch value.Type() {
		case filter.IntType:
			q := bluge.NewNumericRangeQuery(0-math.MaxFloat64, float64(value.Value().(int64)))
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := bluge.NewNumericRangeQuery(0-math.MaxFloat64, value.Value().(float64))
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the bleve mapping
			q := bluge.NewTermRangeQuery("", stripQuotes(value.String()))
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	case filter.CompOpGe:
		switch value.Type() {
		case filter.IntType:
			q := bluge.NewNumericRangeInclusiveQuery(float64(value.Value().(int64)), math.MaxFloat64, true, true)
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := bluge.NewNumericRangeInclusiveQuery(value.Value().(float64), math.MaxFloat64, true, true)
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the mapping
			q := bluge.NewTermRangeInclusiveQuery(stripQuotes(value.String()), "", true, true)
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	case filter.CompOpLe:
		switch value.Type() {
		case filter.IntType:
			q := bluge.NewNumericRangeInclusiveQuery(0-math.MaxFloat64, float64(value.Value().(int64)), true, true)
			q.SetField(field)
			return q, nil
		case filter.FloatType:
			q := bluge.NewNumericRangeInclusiveQuery(0-math.MaxFloat64, value.Value().(float64), true, true)
			q.SetField(field)
			return q, nil
		case filter.StringType:
			// todo: how to handle dates? they don't have a special type so we would need to look
			// at the bleve mapping
			q := bluge.NewTermRangeInclusiveQuery("", stripQuotes(value.String()), true, true)
			q.SetField(field)
			return q, nil
		default:
			return nil, fmt.Errorf("value type %s is not applicable to %s operation", string(value.Type()), string(op))
		}
	default:
		return nil, fmt.Errorf("operation %s was not implemented", string(op))
	}
}

func (j *BlugeQuery) eqFilter(field string, value filter.Value) (bluge.Query, error) {
	if t, ok := mapping.Mapping[field]; ok {
		switch t {
		case mapping.FieldTypeText:
			if value.Type() != filter.StringType {
				return nil, fmt.Errorf("could not compare text field %s with %s", field, value.Type())
			}
			q := bluge.NewMatchPhraseQuery(stripQuotes(value.String()))
			q.SetField(field)
			return q, nil
		case mapping.FieldTypeKeyword:
			if value.Type() != filter.StringType {
				return nil, fmt.Errorf("could not compare keyword field %s with %s", field, value.Type())
			}
			q := bluge.NewTermQuery(stripQuotes(value.String()))
			q.SetField(field)
			return q, nil
		case mapping.FieldTypeNumber:
			switch value.Type() {
			case filter.IntType:
				q := bluge.NewNumericRangeInclusiveQuery(float64(value.Value().(int64)), float64(value.Value().(int64)), true, true)
				q.SetField(field)
				return q, nil
			case filter.FloatType:
				q := bluge.NewNumericRangeInclusiveQuery(value.Value().(float64), value.Value().(float64), true, true)
				q.SetField(field)
				return q, nil
			default:
				return nil, fmt.Errorf("cannot compare number to %s", value.Type())
			}
		case mapping.FieldTypeDate:
			if v, ok := value.Value().(string); ok {
				ts, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return nil, fmt.Errorf("failed to parse %s as date: %s", field, err.Error())
				}
				q := bluge.NewDateRangeQuery(ts, ts)
				q.SetField(field)
				return q, nil
			}
			return nil, fmt.Errorf("non-string value given as date")
		}
	}
	return nil, fmt.Errorf("unknown field %s", field)
}

func stripQuotes(v string) string {
	return strings.Trim(v, `"`)
}
