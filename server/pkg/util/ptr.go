package util

import "time"

func TimeP(t time.Time) *time.Time {
	return &t
}

func Int64P(v int64) *int64 {
	return &v
}
