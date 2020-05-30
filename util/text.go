package util

import "regexp"

func SafeText(text string) string {
	re := regexp.MustCompile("[[:^ascii:]]")
	text = re.ReplaceAllLiteralString(text, "")
	re = regexp.MustCompile("\\s+") // this one necessary
	text = re.ReplaceAllLiteralString(text, " ")
	return text
}
