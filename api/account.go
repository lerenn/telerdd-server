package api

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
)

type Account struct {
	// Infos
	data   *Data
	db     *sql.DB
	logger *log.Log
	// API
	token *Token
}

func newAccount(data *Data, db *sql.DB, logger *log.Log) *Account {
	var a Account
	a.data = data
	a.db = db
	a.logger = logger
	a.token = newToken(data, db, logger)
	return &a
}

func (a *Account) Process(request string, r *http.Request) string {
	base, _ := splitString(request, "/")

	if base == "token" {
		return a.token.Process(r)
	} else {
		return jsonBadURL(r)
	}
}
