package api

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
)

type accountBundle struct {
	// Infos
	data   *data
	db     *sql.DB
	logger *log.Log
	// API
	token *tokenBundle
}

func newAccountBundle(d *data, db *sql.DB, logger *log.Log) *accountBundle {
	var a accountBundle
	a.data = d
	a.db = db
	a.logger = logger
	a.token = newTokenBundle(d, db, logger)
	return &a
}

func (a *accountBundle) Process(request string, r *http.Request) string {
	base, _ := splitString(request, "/")

	if base == "token" {
		return a.token.Process(r)
	}
	return jsonBadURL(r)
}
