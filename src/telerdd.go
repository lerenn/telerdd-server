package main

import (
	"fmt"
	"net/http"

	"github.com/lerenn/telerdd-server/src/api"
)

const NO_DATE = "[###]"

func main() {
	fmt.Println(NO_DATE + " App launched")

	// Prepare config
	conf, err := initConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(NO_DATE + " Configuration loaded")

	// Prepare log
	logger, err := initLog(conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	logger.Print("Log file loaded")

	// Prepare DB
	db, err := initDB(conf)
	if err != nil {
		logger.Print(err.Error())
	}
	logger.Print("Database loaded")

	// Prepare API
	_, err = api.New(conf, db, logger)
	if err != nil {
		logger.Print(err.Error())
	}
	logger.Print("API Server ready")

	// Launch server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Print(err.Error())
	}
}
