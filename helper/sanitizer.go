package helper

import "strings"

func Sanitize(s string) string {
	return strings.Replace(s, "\n", "", -1)
}
