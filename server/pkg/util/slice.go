package util

func InStrings(s string, ss...string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func LastString(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return ss[len(ss)-1]
}

