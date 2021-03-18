package common

import "fmt"

func LimitStmnt(pageSize int32, page int32) string {
	if pageSize < 1 {
		pageSize = 25
	}
	if page < 1 {
		page = 1
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, pageSize*(page-1))
}

