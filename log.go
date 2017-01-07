package main

import (
	"errors"

	libConfig "github.com/lerenn/go-config"
	"github.com/lerenn/log"
	appConfig "github.com/nightwall/nightwall-server/config"
)

func newLog(c *libConfig.Config) (*log.Log, error) {
	// Get log file name
	logFile, err := c.GetString(appConfig.LogSectionToken, appConfig.LogFileToken)
	if err != nil {
		return nil, err
	}

	// Create logger
	logger := log.New()
	if logger.Start(logFile) != nil {
		return logger, errors.New("Error when trying to create a new log file")
	}

	return logger, nil
}
