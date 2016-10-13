package api

import (
	"database/sql"

	"github.com/lerenn/log"
)

type bundle struct {
	data   *data
	db     *sql.DB
	logger *log.Log
}

func newBundle(d *data, db *sql.DB, logger *log.Log) bundle {
	return bundle{d, db, logger}
}
