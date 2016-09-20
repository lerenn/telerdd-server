package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	config "github.com/lerenn/go-config"
)

func initDB(c *config.Config) (*sql.DB, error) {
	var user, password, addr, port, name string
	var err error

	// Get params
	if user, err = c.GetString(DB_SECTION_TOKEN, DB_USER_TOKEN); err != nil {
		return nil, err
	} else if password, err = c.GetString(DB_SECTION_TOKEN, DB_PASSWORD_TOKEN); err != nil {
		return nil, err
	} else if addr, err = c.GetString(DB_SECTION_TOKEN, DB_ADDR_TOKEN); err != nil {
		return nil, err
	} else if port, err = c.GetString(DB_SECTION_TOKEN, DB_PORT_TOKEN); err != nil {
		return nil, err
	} else if name, err = c.GetString(DB_SECTION_TOKEN, DB_NAME_TOKEN); err != nil {
		return nil, err
	}

	// Open database
	dataSourceName := user + ":" + password + "@tcp(" + addr + ":" + port + ")/" + name // user:password@tcp(addr:port)/db
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
