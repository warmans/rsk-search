package util

import (
	"fmt"
	"strings"
)

func CreatePlaceholdersForStrings(ss []string) (string, []interface{}) {
	ph := make([]string, len(ss))
	params := make([]interface{}, len(ss))
	for k, v := range ss {
		ph[k] = fmt.Sprintf("$%d", k+1)
		params[k] = v
	}
	return strings.Join(ph, ", "), params
}
