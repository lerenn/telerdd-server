package config

import (
	"errors"
	"fmt"
	"strings"
)

type Section struct {
	name   string
	values map[string]string
}

// Create a new section
// param: name Name of the section
// return: Pointer to a new section
func newSection(name string) *Section {
	var sec Section
	sec.name = name
	sec.values = make(map[string]string)
	return &sec
}

// Add a token/value
// param: token Token to add
// param: value Value to add
func (s *Section) Add(token, value string) {
	s.values[token] = value
}

// Get section name
// return: section name
func (s *Section) Name() string {
	return s.name
}

// Get string from section
// param: token Token searched
// param: defaultValue defaultValue displayed if there is an error
// return: Value corresponding to token and error state
func (s *Section) GetString(token string) (string, error) {
	value, presence := s.values[token]
	if presence {
		return trimValue(value), nil
	} else {
		err := fmt.Sprintf("No %q token found", token)
		return "", errors.New(err)
	}
}

// Trim value with useless characters
// param: value Value to trim
// return: value without start/end useless characters
func trimValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "\"")
	value = strings.TrimSuffix(value, "\"")
	value = strings.TrimPrefix(value, "'")
	value = strings.TrimSuffix(value, "'")
	return value
}
