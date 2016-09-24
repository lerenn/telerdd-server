package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/api/account"
	"github.com/lerenn/telerdd-server/src/api/messages"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type API struct {
	// Infos
	data             *data.Data
	db               *sql.DB
	logger           *log.Log
	authorizedOrigin string
	// API
	account  *account.Account
	messages *messages.Messages
}

func New(data *data.Data, db *sql.DB, logger *log.Log, authorizedOrigin string) *API {
	var a API
	a.account = account.New(data, db, logger)
	a.data = data
	a.db = db
	a.logger = logger
	a.authorizedOrigin = authorizedOrigin
	a.messages = messages.New(data, db, logger)
	return &a
}

func (a *API) Process(w http.ResponseWriter, r *http.Request) {
	base, extent := tools.Split(r.URL.Path[1:], "/")

	// Authorize origin
	w.Header().Set("Access-Control-Allow-Origin", a.authorizedOrigin)

	var response string
	if base == "messages" {
		response = a.messages.Process(extent, r)
	} else if base == "account" {
		response = a.account.Process(extent, r)
	} else {
		response = tools.JSONBadURL(r)
	}

	fmt.Fprintf(w, response)
}
