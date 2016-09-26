package main

import (
	"errors"

	config "github.com/lerenn/go-config"
	"github.com/lerenn/log"
	cst "github.com/lerenn/telerdd-server/src/constants"
)

func initLog(c *config.Config) (*log.Log, error) {
	// Get log file name
	logFile, err := c.GetString(cst.LOG_SECTION_TOKEN, cst.LOG_FILE_TOKEN)
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
