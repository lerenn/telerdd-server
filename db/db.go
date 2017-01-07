package db

import (
	"database/sql"

	libConfig "github.com/lerenn/go-config"
	appConfig "github.com/nightwall/nightwall-server/config"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

// New database for application
func New(c *libConfig.Config) (*sql.DB, error) {
	var user, password, addr, port, name string
	var err error

	// Get params
	if user, err = c.GetString(appConfig.DbSectionToken, appConfig.DbUserToken); err != nil {
		return nil, err
	} else if password, err = c.GetString(appConfig.DbSectionToken, appConfig.DbPasswordToken); err != nil {
		return nil, err
	} else if addr, err = c.GetString(appConfig.DbSectionToken, appConfig.DbAddrToken); err != nil {
		return nil, err
	} else if port, err = c.GetString(appConfig.DbSectionToken, appConfig.DbPortToken); err != nil {
		return nil, err
	} else if name, err = c.GetString(appConfig.DbSectionToken, appConfig.DbNameToken); err != nil {
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
