package common

import (
	"fmt"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/filter/psql"
)

type Pager interface {
	GetPage() int32
	GetPageSize() int32
}

type Sorter interface {
	GetSortField() string
	GetSortDirection() string
}

type Filterer interface {
	GetFilter() string
}

type SortDirection string

const SortAsc SortDirection = "ASC"
const SortDesc SortDirection = "DESC"

type Sorting struct {
	Field     string
	Direction SortDirection
}

func (p *Sorting) Stmnt(fieldMap map[string]string) (string, error) {
	if p == nil || p.Field == "" {
		return "", nil
	}
	f, ok := fieldMap[p.Field]
	if !ok {
		return "", fmt.Errorf("%s is not a sortable field", p.Field)
	}
	if p.Direction != SortAsc && p.Direction != SortDesc && p.Direction != "" {
		return "", fmt.Errorf("%s is not a sort direction", p.Direction)
	}
	return fmt.Sprintf("ORDER BY %s %s", f, string(p.Direction)), nil
}

func NewDefaultPaging() *Paging {
	return &Paging{
		Page:     0,
		PageSize: 25,
	}
}

type Paging struct {
	Page     int32
	PageSize int32
}

func (p *Paging) Stmnt() string {
	if p == nil {
		return ""
	}
	return LimitStmnt(p.PageSize, p.Page)
}

type QueryOpt func(q *QueryModifier)

func WithPaging(pageSize, page int32) QueryOpt {
	return func(q *QueryModifier) {
		q.Paging = &Paging{Page: page, PageSize: pageSize}
	}
}

func WithSorting(field string, direction SortDirection) QueryOpt {
	return func(q *QueryModifier) {
		q.Sorting = &Sorting{Field: field, Direction: direction}
	}
}

func WithFilter(f filter.Filter) QueryOpt {
	return func(q *QueryModifier) {
		q.Filter = f
	}
}

func WithDefaultSorting(field string, direction SortDirection) QueryOpt {
	return func(q *QueryModifier) {
		q.defaultSorting = &Sorting{Field: field, Direction: direction}
	}
}

func Q(opts ...QueryOpt) *QueryModifier {
	q := &QueryModifier{}
	for _, v := range opts {
		v(q)
	}
	return q
}

// QueryModifier - adds filtering, sorting paging to queries.
type QueryModifier struct {
	Paging  *Paging
	Sorting *Sorting
	Filter  filter.Filter

	// if no sorting is passed allow a default to be given.
	defaultSorting *Sorting
}

func (q *QueryModifier) Apply(opt QueryOpt) *QueryModifier {
	opt(q)
	return q
}

func (q *QueryModifier) ToSQL(fieldMap map[string]string, withWHERE bool) (where string, params []interface{}, order string, paging string, err error) {
	if q == nil {
		return "", []interface{}{}, "", "", nil
	}
	if q.Filter != nil {
		where, params, err = psql.FilterToQuery(q.Filter, fieldMap)
		if err != nil {
			return
		}
		if withWHERE {
			where = fmt.Sprintf("WHERE %s", where)
		}
	}
	if q.Sorting != nil {
		order, err = q.Sorting.Stmnt(fieldMap)
		if err != nil {
			return
		}
	} else {
		if q.defaultSorting != nil {
			order, err = q.defaultSorting.Stmnt(fieldMap)
			if err != nil {
				return
			}
		}
	}
	if q.Paging != nil {
		paging = q.Paging.Stmnt()
	}
	return
}

func LimitStmnt(pageSize int32, page int32) string {
	if pageSize < 1 {
		pageSize = 25
	}
	if page < 1 {
		page = 1
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, pageSize*(page-1))
}
