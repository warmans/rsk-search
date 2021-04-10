package util

import "time"

func TimeP(t time.Time) *time.Time {
	return &t
}

func Int64P(v int64) *int64 {
	return &v
}

func PString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func PFloat32(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}

func PTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
