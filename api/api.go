package api

import (
	"database/sql"
	"fmt"
	"net/http"

	config "github.com/lerenn/go-config"
	"github.com/lerenn/log"
)

type API struct {
	data     *Data
	account  *Account
	messages *Messages
}

func New(c *config.Config, db *sql.DB, logger *log.Log) (*API, error) {
	var a API
	var err error

	// Create components
	if a.data, err = newData(c); err != nil {
		return nil, err
	}
	a.account = newAccount(a.data, db, logger)
	a.messages = newMessages(a.data, db, logger)

	// Set callback
	http.HandleFunc("/", a.Process)

	return &a, nil
}

func (a *API) Process(w http.ResponseWriter, r *http.Request) {
	base, extent := splitString(r.URL.Path[1:], "/")

	// Set header
	a.setHeader(w)

	// Process response
	var response string
	if base == "messages" {
		response = a.messages.Process(extent, r)
	} else if base == "account" {
		response = a.account.Process(extent, r)
	} else {
		response = jsonBadURL(r)
	}

	fmt.Fprintf(w, response)
}

func (a *API) setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", a.data.AuthorizedOrigin())
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
}
