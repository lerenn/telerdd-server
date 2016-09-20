package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Redifining some functions for mocking
var readFile = ioutil.ReadFile

// Config structure
type Config struct {
	sections []*Section
}

// Create a new configuration structure
// return: Configuration structure
func New() *Config {
	var conf Config
	conf.sections = make([]*Section, 0)
	return &conf
}

// Read a config from a file
// param: Path to the file
// return: nil if successful, false otherwise
func (c *Config) Read(filePath string) error {
	// Reading configuration file
	content, err := readFile(filePath)
	if err != nil {
		return errors.New("Error when reading configuration file.")
	}
	lines := strings.Split(string(content), "\n")

	// Process lines
	var section *Section = nil
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if newSec := newSectionFromLine(line); newSec != nil { // If the line is a section
			section = newSec
			c.sections = append(c.sections, section)
		} else if newTok, newVal := getTokenValue(line); newTok != "" { // If the line is a token/value
			section.Add(newTok, newVal)
		}
	}
	return nil
}

// Create a new section if the line contains a section title
// param: line Line to examinate
// return: a pointer to a new section if it exists, nil otherwise
func newSectionFromLine(line string) *Section {
	begin := strings.Index(line, "[")
	end := strings.LastIndex(line, "]")

	if begin != -1 && end != -1 && begin != end {
		sec := newSection(line[begin+1 : end])
		return sec
	}

	return nil
}

// Get token and value from a line, if possible
// param: line Line to examinate
// return: token and its value
func getTokenValue(line string) (string, string) {
	tokenValue := strings.SplitN(line, "=", 2)
	switch len(tokenValue) {
	case 2:
		return tokenValue[0], tokenValue[1]
	default:
		return tokenValue[0], ""
	}
}

// Get string from configuration
// param: section Section where is the token
// param: token Token searched
// return: Value corresponding to token and error state
func (c *Config) GetString(section, token string) (string, error) {
	for _, s := range c.sections {
		if s.Name() == section {
			return s.GetString(token)
		}
	}

	err := fmt.Sprintf("No %q section found", section)
	return "", errors.New(err)
}

// Get int from configuration
// param: token Token searched
// return: Value corresponding to token and error state
func (c *Config) GetInt(section, token string) (int, error) {
		valueStr, err := c.GetString(section, token)
		if err != nil {
			return 0, err
		}

		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return 0, err
		}

		return value, nil
}
