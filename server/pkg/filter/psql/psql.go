package psql

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/filter"
	"strings"
)

func FilterToQuery(f filter.Filter, filterMapping map[string]string) (string, []interface{}, error) {
	v := &visitor{filterMapping: filterMapping, params: newParams()}
	err := f.Accept(v)
	return v.sql, v.params.par, err
}

func filterAppendQuery(f filter.Filter, filterMapping map[string]string, params *params) (string, []interface{}, error) {
	v := &visitor{filterMapping: filterMapping, params: params}
	err := f.Accept(v)
	return v.sql, v.params.par, err
}

type visitor struct {
	filterMapping map[string]string

	sql    string
	params *params
}

func (j *visitor) VisitCompFilter(f *filter.CompFilter) (filter.Visitor, error) {
	qb := newQueryBuilder(j.filterMapping)
	var err error
	j.sql, err = qb.compExpr(f, j.params)
	return nil, err
}

func (j *visitor) VisitBoolFilter(f *filter.BoolFilter) (filter.Visitor, error) {
	qb := newQueryBuilder(j.filterMapping)
	var err error
	j.sql, err = qb.boolExpr(f, j.params)
	return nil, err
}

func newQueryBuilder(filterMapping map[string]string) *queryBuilder {
	return &queryBuilder{
		filterToSelect: filterMapping,
	}
}

type queryBuilder struct {
	sql            string
	filterToSelect map[string]string
}

func (j *queryBuilder) boolExpr(f *filter.BoolFilter, p *params) (string, error) {

	lhs, _, err := filterAppendQuery(f.LHS, j.filterToSelect, p)
	if err != nil {
		return "", err
	}

	rhs, _, err := filterAppendQuery(f.RHS, j.filterToSelect, p)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s) %s (%s)", lhs, f.Op, rhs), nil
}

func (j *queryBuilder) compExpr(f *filter.CompFilter, p *params) (string, error) {
	col, ok := j.filterToSelect[f.Field]
	if !ok {
		return "", errors.Errorf("field is not filterable: '%s'", f.Field)
	}

	switch comp := f.Op; comp {
	case filter.CompOpEq:
		if f.Value.IsNull() {
			return fmt.Sprintf(`%s IS NULL`, col), nil
		}
		return fmt.Sprintf(`%s = %s`, col, p.next(f.Value.Value())), nil
	case filter.CompOpNeq:
		if f.Value.IsNull() {
			return fmt.Sprintf(`%s IS NOT NULL`, col), nil
		}
		return fmt.Sprintf(`%s != %s`, col, p.next(f.Value.Value())), nil
	case filter.CompOpGt:
		return fmt.Sprintf(`%s > %s`, col, p.next(f.Value.Value())), nil
	case filter.CompOpGe:
		return fmt.Sprintf(`%s >= %s`, col, p.next(f.Value.Value())), nil
	case filter.CompOpLt:
		return fmt.Sprintf(`%s < %s`, col, p.next(f.Value.Value())), nil
	case filter.CompOpLe:
		return fmt.Sprintf(`%s <= %s`, col, p.next(f.Value.Value())), nil
	case filter.CompOpLike:
		return fmt.Sprintf(`%s LIKE %s`, col, p.next(fmt.Sprintf("%%%s%%", strings.Trim(f.Value.String(), `"`)))), nil
	case filter.CompOpFuzzyLike:
		// could make this better if proper fuzzy search is needed.
		return fmt.Sprintf(`%s ILIKE %s`, col, p.next(fmt.Sprintf("%%%s%%", strings.Trim(f.Value.String(), `"`)))), nil
	}
	return "", fmt.Errorf("unknown operator: %s", string(f.Op))
}

func newParams() *params {
	return &params{par: []interface{}{}}
}

type params struct {
	par []interface{}
}

func (p *params) next(v interface{}) string {
	p.par = append(p.par, v)
	return fmt.Sprintf("$%d", len(p.par))
}

func (p *params) add(ps ...interface{}) {
	for _, b := range ps {
		p.par = append(p.par, b)
	}
}
