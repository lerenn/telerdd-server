package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	config "github.com/lerenn/go-config"
	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/api"
	"github.com/lerenn/telerdd-server/src/data"

	_ "github.com/go-sql-driver/mysql"
)

type TeleRDD struct {
	api    *api.API
	conf   *config.Config
	logger *log.Log
	db     *sql.DB
	data   *data.Data
}

func (t *TeleRDD) Init() error {
	// Prepare config
	if err := t.initConfig(); err != nil {
		return err
	}
	fmt.Println(NO_DATE + " Configuration loaded")

	// Prepare log
	if err := t.initLog(); err != nil {
		return err
	}
	t.logger.Print("Log file loaded")

	// Prepare API infos
	if err := t.initData(); err != nil {
		return err
	}
	t.logger.Print("API info loaded from conf")

	// Prepare DB
	if err := t.initDB(); err != nil {
		return err
	}
	t.logger.Print("Database loaded")

	// Prepare API
	if err := t.initAPI(); err != nil {
		return err
	}
	t.logger.Print("API Server ready")

	t.logger.Print("Initialization complete")
	return nil
}

func (t *TeleRDD) Start() error {
	return http.ListenAndServe(":8080", nil)
}

func (t *TeleRDD) CloseDB() {
	t.db.Close()
}

// Private methods
////////////////////////////////////////////////////////////////////////////////

func (t *TeleRDD) initAPI() error {
	// Get authorized URL for client
	authorizedOrigin, err := t.conf.GetString(CLIENT_SECTION_TOKEN, CLIENT_AUTHORIZED_ORIGIN_TOKEN)
	if err != nil {
		return err
	}

	// Create api
	t.api = api.New(t.data, t.db, t.logger, authorizedOrigin)

	// Set callback
	http.HandleFunc("/", t.api.Process)

	return nil
}

func (t *TeleRDD) initConfig() error {
	// Read the conf file
	t.conf = config.New()
	if err := t.conf.Read(CONFIG_FILE); err != nil {
		return errors.New("Error when reading config file.")
	}

	return nil
}

func (t *TeleRDD) initDB() error {
	var user, password, addr, port, name string
	var err error
	c := t.conf

	// Get params
	if user, err = c.GetString(DB_SECTION_TOKEN, DB_USER_TOKEN); err != nil {
		return err
	} else if password, err = c.GetString(DB_SECTION_TOKEN, DB_PASSWORD_TOKEN); err != nil {
		return err
	} else if addr, err = c.GetString(DB_SECTION_TOKEN, DB_ADDR_TOKEN); err != nil {
		return err
	} else if port, err = c.GetString(DB_SECTION_TOKEN, DB_PORT_TOKEN); err != nil {
		return err
	} else if name, err = c.GetString(DB_SECTION_TOKEN, DB_NAME_TOKEN); err != nil {
		return err
	}

	// Open database
	dataSourceName := user + ":" + password + "@tcp(" + addr + ":" + port + ")/" + name // user:password@tcp(addr:port)/db
	if t.db, err = sql.Open("mysql", dataSourceName); err != nil {
		return err
	}

	return nil
}

func (t *TeleRDD) initLog() error {
	// Get log file name
	logFile, err := t.conf.GetString(LOG_SECTION_TOKEN, LOG_FILE_TOKEN)
	if err != nil {
		return err
	}

	// Create logger
	t.logger = log.New()
	if t.logger.Start(logFile) != nil {
		return errors.New("Error when trying to create a new log file")
	}

	return nil
}

func (t *TeleRDD) initData() error {
	msgLimit, err := t.conf.GetInt(MESSAGES_SECTION_TOKEN, MESSAGES_LIMIT_TOKEN)
	if err != nil {
		return err
	}

	t.data = data.New(msgLimit)
	return nil
}
