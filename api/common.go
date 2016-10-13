package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tomasen/realip"
)

// JSON
////////////////////////////////////////////////////////////////////////////////

func jsonResponseOk() string {
	return "{\"response\": \"OK\"}"
}

func jsonError(txt string) string {
	return "{\"error\":\"" + replaceBadChar(txt) + "\"}"
}

func jsonBadURL(r *http.Request) string {
	url := r.URL.Path[1:]
	return jsonError("Bad URL: " + url)
}

// Request operations
////////////////////////////////////////////////////////////////////////////////

func getRequestIP(r *http.Request) string {
	return realip.RealIP(r)
}

func getRequestInt(r *http.Request, name string) (int, bool, error) {
	// Get string
	nbrStr := r.FormValue(name)

	// If string is empty, then it is not passed in parameter
	if nbrStr == "" {
		return 0, false, nil
	}

	// Try to change string to nbr
	nbr, err := strconv.Atoi(nbrStr)
	if err != nil {
		errStr := fmt.Sprintf("Invalid number : %q", nbrStr)
		return 0, true, errors.New(errStr)
	}

	// Return number
	return nbr, true, nil
}

// SQL Operations
////////////////////////////////////////////////////////////////////////////////

const sqlDatetimeForm = "2006-01-02 15:04:05"

func sqlTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func sqlFormatDateTime(orig string) (string, error) {
	t, err := sqlParseTime(orig)
	if err != nil {
		return "", errors.New("Error when parsing date")
	}
	return t.Format("02/01/2006 - 15:04:05"), nil
}

func sqlParseTime(t string) (time.Time, error) {
	return time.Parse(sqlDatetimeForm, t)
}

// Strings
////////////////////////////////////////////////////////////////////////////////

func replaceBadChar(s string) string {
	s = strings.Replace(s, "\r", "", -1)     // Replace \r by nothing
	s = strings.Replace(s, "\n", "\\n", -1)  // Replace \n by newline character
	s = strings.Replace(s, "\"", "\\\"", -1) // Replace " by \"
	return s
}

func splitString(line, separator string) (string, string) {
	index := strings.Index(line, separator)
	if index < 0 {
		return line, ""
	}
	return line[:index], line[index+len(separator):]
}
