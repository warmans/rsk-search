package util

import "time"

const ShortDateFormat = "2006-01-02"
const SQLDateFormat = "2006-01-02 15:04:05"

func ShortDate(t time.Time) string {
	return t.Format(ShortDateFormat)
}

func SqlDate(t time.Time) string {
	return t.Format(SQLDateFormat)
}
