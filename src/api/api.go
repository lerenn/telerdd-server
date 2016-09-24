package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/api/v1"
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
	v1 *v1.V1
}

func New(data *data.Data, db *sql.DB, logger *log.Log, authorizedOrigin string) *API {
	var a API
	a.data = data
	a.db = db
	a.logger = logger
	a.authorizedOrigin = authorizedOrigin
	a.v1 = v1.New(data, db, logger)
	return &a
}

func (a *API) Process(w http.ResponseWriter, r *http.Request) {
	base, extent := tools.Split(r.URL.Path[1:], "/")

	// Set header
	a.setHeader(w)

	var response string
	if base == "v1" {
		response = a.v1.Process(extent, r)
	} else {
		response = tools.JSONBadURL(r)
	}

	fmt.Fprintf(w, response)
}

func (a *API) setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", a.authorizedOrigin)
	w.Header().Set("Content-Type", "application/json")
}
