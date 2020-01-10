package ocrmutils

import (
	"regexp"
	"strings"
)

var (
	digitRegExp *regexp.Regexp
)

func init() {
	digitRegExp = regexp.MustCompile("[0-9]+")
}

func TrimAllSpacesInString(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.Join(strings.Fields(str), "")
}

func ReplaceCharInStringAtIndex(str string, char rune, index int) string {
	if len(str) == 0 {
		return ""
	}
	out := []rune(str)
	out[index] = char
	return string(out)
}

func GetDigitsFromString(value string) string {
	return strings.Join(digitRegExp.FindAllString(value, -1), "")
}
