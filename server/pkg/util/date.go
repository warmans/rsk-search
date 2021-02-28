package util

import "time"

const ShortDateFormat = "2006-01-02"


func ShortDate(t time.Time) string {
	return t.Format(ShortDateFormat)
}
