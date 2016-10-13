package api

import (
	"database/sql"
	"fmt"
	"net/http"

	config "github.com/lerenn/go-config"
	"github.com/lerenn/log"
)

// API instance
type API struct {
	data     *data
	account  accountBundle
	messages messagesBundle
}

// New API instance
func New(c *config.Config, db *sql.DB, logger *log.Log) (*API, error) {
	var a API
	var err error

	// Create components
	if a.data, err = newData(c); err != nil {
		return nil, err
	}
	b := newBundle(a.data, db, logger)
	a.account = newAccountBundle(b)
	a.messages = newMessagesBundle(b)

	// Set callback
	http.HandleFunc("/", a.Process)

	return &a, nil
}

// Process HTTP Request
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
