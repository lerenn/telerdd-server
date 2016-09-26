package account

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
)

type Account struct {
	// Infos
	data   *common.Data
	db     *sql.DB
	logger *log.Log
	// API
	token *Token
}

func New(data *common.Data, db *sql.DB, logger *log.Log) *Account {
	var a Account
	a.data = data
	a.db = db
	a.logger = logger
	a.token = NewToken(data, db, logger)
	return &a
}

func (a *Account) Process(request string, r *http.Request) string {
	base, _ := common.Split(request, "/")

	if base == "token" {
		return a.token.Process(r)
	} else {
		return common.JSONBadURL(r)
	}
}
