package utils

import (
	"regexp"
)

// RegexpToString converts *regexp.Regexp instance to string
func RegexpToString(r *regexp.Regexp) string {
	if r != nil {
		return r.String()
	}
	return ""
}
