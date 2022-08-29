package helper

import "strings"

func Sanitize(s string) string {
	return strings.Replace(s, "\n", "", -1)
}

func SanitizeStringSlice(ss []string) {
	for i, x := range ss {
		ss[i] = Sanitize(x)
	}
}
