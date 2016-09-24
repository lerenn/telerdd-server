package v1

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/api/v1/account"
	"github.com/lerenn/telerdd-server/src/api/v1/messages"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type V1 struct {
	// Infos
	data   *data.Data
	db     *sql.DB
	logger *log.Log
	// API
	account  *account.Account
	messages *messages.Messages
}

func New(data *data.Data, db *sql.DB, logger *log.Log) *V1 {
	var v V1
	v.data = data
	v.db = db
	v.logger = logger
	v.account = account.New(data, db, logger)
	v.messages = messages.New(data, db, logger)
	return &v
}

func (v *V1) Process(request string, r *http.Request) string {
	base, extent := tools.Split(request, "/")

	if base == "messages" {
		return v.messages.Process(extent, r)
	} else if base == "account" {
		return v.account.Process(extent, r)
	} else {
		return tools.JSONBadURL(r)
	}
}
