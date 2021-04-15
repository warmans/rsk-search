package common

import "fmt"

type SortDirection string

const SortAsc SortDirection = "ASC"
const SortDESC SortDirection = "DESC"

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
	if p.Direction != SortAsc && p.Direction != SortDESC && p.Direction != "" {
		return "", fmt.Errorf("%s is not a sort direction", p.Direction)
	}
	return fmt.Sprintf("ORDER BY %s %s", f, string(p.Direction)), nil
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

func LimitStmnt(pageSize int32, page int32) string {
	if pageSize < 1 {
		pageSize = 25
	}
	if page < 1 {
		page = 1
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, pageSize*(page-1))
}
