package messages

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/telerdd/src/api/account"
	"github.com/lerenn/telerdd/src/tools"
)

func status(db *sql.DB, r *http.Request) string {
	// Check permission
	errAuth := account.Authorized(db, r, 1)
	if errAuth != nil {
		return tools.JSONError(errAuth.Error())
	}

	// Get id
	id, err := tools.GetIntFromRequest(r, "id")
	if err != nil {
		return tools.JSONError("Invalid number for id")
	} else if id < 0 {
		id = 0
	}

	// Get status
	status := r.FormValue("status")
	if status != "accepted" && status != "refused" {
		return tools.JSONError("Invalid status")
	}

	// Prepare request
	stmt, errPrep := db.Prepare("UPDATE messages SET status=? WHERE id=?")
	if errPrep != nil {
		return tools.JSONError(errPrep.Error())
	}

	_, errExec := stmt.Exec(status, id)
	if errExec != nil {
		return tools.JSONError(errExec.Error())
	}

	return tools.JSONResponseOk()
}
