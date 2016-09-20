package api

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd/src/api/account"
	"github.com/lerenn/telerdd/src/api/messages"
	"github.com/lerenn/telerdd/src/data"
	"github.com/lerenn/telerdd/src/tools"
)

const Prefix = "/api/"

func Process(infos *data.Data, db *sql.DB, l *log.Log, w http.ResponseWriter, r *http.Request, request string) string {
	base, extent := tools.Split(request, "/")

	if base == "messages" {
		return messages.Process(infos, db, l, extent, r)
	} else if base == "account" {
		return account.Process(infos, db, l, extent, r)
	} else {
		return tools.JSONBadURL(r)
	}
}
