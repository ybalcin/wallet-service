package utility

import "strings"

func StrLength(str string) int {
	str = strings.TrimSpace(str)
	return len([]rune(str))
}

func IsStrEmpty(str string) bool {
	return StrLength(str) <= 0
}
