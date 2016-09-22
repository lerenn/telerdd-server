package account

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

func Process(infos *data.Data, db *sql.DB, l *log.Log, request string, r *http.Request) string {
	base, _ := tools.Split(request, "/")

	if base == "connect" {
		return connect(infos, db, l, r)
	} else {
		return tools.JSONBadURL(r)
	}
}
