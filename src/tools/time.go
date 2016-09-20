package tools

import (
	"errors"
	"time"
)

const SQL_DATETIME_FORM = "2006-01-02 15:04:05"

func SQLTimeNow() string {
	return time.Now().Format(SQL_DATETIME_FORM)
}

func SQLFormatDateTime(orig string) (string, error) {
	t, err := SQLParseTime(orig)
	if err != nil {
		return "", errors.New("Error when parsing date")
	}
	return t.Format("02/01/2006 - 15:04:05"), nil
}

func SQLParseTime(t string) (time.Time, error) {
	return time.Parse(SQL_DATETIME_FORM, t)
}
