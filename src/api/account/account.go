package account

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type Account struct {
	// Infos
	data   *data.Data
	db     *sql.DB
	logger *log.Log
	// API
	connect *Connect
}

func New(data *data.Data, db *sql.DB, logger *log.Log) *Account {
	var a Account
	a.data = data
	a.db = db
	a.logger = logger
	a.connect = NewConnect(data, db, logger)
	return &a
}

func (a *Account) Process(request string, r *http.Request) string {
	base, _ := tools.Split(request, "/")

	if base == "connect" {
		return a.connect.Process(r)
	} else {
		return tools.JSONBadURL(r)
	}
}
