package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	config "github.com/lerenn/go-config"
	cst "github.com/lerenn/telerdd-server/constants"
)

func New(c *config.Config) (*sql.DB, error) {
	var user, password, addr, port, name string
	var err error

	// Get params
	if user, err = c.GetString(cst.DB_SECTION_TOKEN, cst.DB_USER_TOKEN); err != nil {
		return nil, err
	} else if password, err = c.GetString(cst.DB_SECTION_TOKEN, cst.DB_PASSWORD_TOKEN); err != nil {
		return nil, err
	} else if addr, err = c.GetString(cst.DB_SECTION_TOKEN, cst.DB_ADDR_TOKEN); err != nil {
		return nil, err
	} else if port, err = c.GetString(cst.DB_SECTION_TOKEN, cst.DB_PORT_TOKEN); err != nil {
		return nil, err
	} else if name, err = c.GetString(cst.DB_SECTION_TOKEN, cst.DB_NAME_TOKEN); err != nil {
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
