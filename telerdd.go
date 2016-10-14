package main

import (
	"fmt"
	"net/http"

	"github.com/lerenn/telerdd-server/api"
	"github.com/lerenn/telerdd-server/config"
	"github.com/lerenn/telerdd-server/db"
)

const noDate = "[###]"

func main() {
	fmt.Println(noDate + " App launched")

	// Prepare config
	conf, err := config.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(noDate + " Configuration loaded")

	// Prepare log
	logger, err := newLog(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logger.Print("Log file loaded")

	// Prepare DB
	db, err := db.New(conf)
	if err != nil {
		logger.Print(err.Error())
		return
	}
	logger.Print("Database loaded")

	// Prepare API
	_, err = api.New(conf, db, logger)
	if err != nil {
		logger.Print(err.Error())
		return
	}
	logger.Print("API Server ready")

	// Launch server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Print(err.Error())
		return
	}
}
