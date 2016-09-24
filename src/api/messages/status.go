package messages

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/api/account"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type Status struct {
	// Infos
	data   *data.Data
	db     *sql.DB
	logger *log.Log
}

func NewStatus(data *data.Data, db *sql.DB, logger *log.Log) *Status {
	var s Status
	s.data = data
	s.db = db
	s.logger = logger
	return &s
}

func (s *Status) Process(r *http.Request) string {
	// Check permission
	errAuth := account.Authorized(s.db, r, 1)
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
	stmt, errPrep := s.db.Prepare("UPDATE messages SET status=? WHERE id=?")
	if errPrep != nil {
		return tools.JSONError(errPrep.Error())
	}

	_, errExec := stmt.Exec(status, id)
	if errExec != nil {
		return tools.JSONError(errExec.Error())
	}

	return tools.JSONResponseOk()
}
