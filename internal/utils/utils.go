package utils

import (
	"strings"
)

// Repeat to repeat string.
func Repeat(str string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(str, n)
}

// PadLeft to pad string to the left.
func PadLeft(str string, l int, p string) string {
	return Repeat(p, l-len([]rune(str))) + str
}

// PadRight to pad string to the right.
func PadRight(str string, l int, p string) string {
	return str + Repeat(p, l-len([]rune(str)))
}

// Ellipsis to truncate string.
func Ellipsis(str string, length int) string {
	r := []rune(str)
	if len(r) > length {
		return string(r[:length-3]) + "..."
	}
	return str
}
