package tools

import (
	"strings"
)

func ReplaceBadCharacters(s string) string {
	s = strings.Replace(s, "\r", "", -1)     // Replace \r by nothing
	s = strings.Replace(s, "\n", "\\n", -1)  // Replace \n by newline character
	s = strings.Replace(s, "\"", "\\\"", -1) // Replace " by \"
	return s
}

func Split(line, separator string) (string, string) {
	index := strings.Index(line, separator)
	if index < 0 {
		return line, ""
	} else {
		return line[:index], line[index+1:]
	}
}
