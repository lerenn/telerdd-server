package messages

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

func Process(infos *data.Data, db *sql.DB, l *log.Log, request string, r *http.Request) string {
	base, _ := tools.Split(request, "/")

	if base == "next" {
		return next(db, r)
	} else if base == "previous" {
		return previous(db, r)
	} else if base == "message" {
		return message(db, r)
	} else if base == "id" {
		return id(db, r)
	} else if base == "send" {
		return send(infos, db, l, r)
	} else if base == "status" {
		return status(db, r)
	} else {
		return tools.JSONBadURL(r)
	}
}
