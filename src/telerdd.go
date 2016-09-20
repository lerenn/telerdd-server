package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	config "github.com/lerenn/go-config"
	"github.com/lerenn/log"
	"github.com/lerenn/telerdd/src/api"
	"github.com/lerenn/telerdd/src/data"
	"github.com/lerenn/telerdd/src/tools"

	_ "github.com/go-sql-driver/mysql"
)

type TeleRDD struct {
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

	// Prepare server
	t.initHTTPServer()
	t.logger.Print("HTTP Server ready")

	t.logger.Print("Initialization complete")
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

func (t *TeleRDD) initHTTPServer() {
	http.Handle(WEBCLIENT_PREFIX, http.StripPrefix(WEBCLIENT_PREFIX, http.FileServer(http.Dir("webclient"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		func(w http.ResponseWriter, r *http.Request, logger *log.Log, db *sql.DB) {
			base, extent := tools.Split(r.URL.Path[1:], "/")

			if base == "api" {
				fmt.Fprintf(w, api.Process(t.data, db, logger, w, r, extent))
			} else {
				http.Redirect(w, r, WEBCLIENT_PREFIX+"home.html", 301)
			}
		}(w, r, t.logger, t.db)
	})
}

func (t *TeleRDD) initData() error {
	msgLimit, err := t.conf.GetInt(MESSAGES_SECTION_TOKEN, MESSAGES_LIMIT_TOKEN)
	if err != nil {
		return err
	}

	t.data = data.New(msgLimit)
	return nil
}

func (t *TeleRDD) Start() error {
	return http.ListenAndServe(":8080", nil)
}

func (t *TeleRDD) CloseDB() {
	t.db.Close()
}
