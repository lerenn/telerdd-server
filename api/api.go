package api

import (
	"database/sql"
	"fmt"
	"net/http"

	config "github.com/lerenn/go-config"
	"github.com/lerenn/log"
	cst "github.com/lerenn/telerdd-server/constants"
)

type API struct {
	data     *Data
	account  *Account
	messages *Messages
}

func New(c *config.Config, db *sql.DB, logger *log.Log) (*API, error) {
	// Get msg limit
	msgLimit, err := c.GetInt(cst.MESSAGES_SECTION_TOKEN, cst.MESSAGES_LIMIT_TOKEN)
	if err != nil {
		return nil, err
	}

	// Get authorized URL for client
	authorizedOrigin, err := c.GetString(cst.CLIENT_SECTION_TOKEN, cst.CLIENT_AUTHORIZED_ORIGIN_TOKEN)
	if err != nil {
		return nil, err
	}

	// Create api data
	data := newData()
	data.SetMsgLimit(msgLimit)
	data.SetAuthorizedOrigin(authorizedOrigin)

	// Create API
	var a API
	a.data = data
	a.account = newAccount(data, db, logger)
	a.messages = newMessages(data, db, logger)

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
