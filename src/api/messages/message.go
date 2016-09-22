package messages

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/telerdd-server/src/tools"
)

func message(db *sql.DB, r *http.Request) string {
	// Get id
	id, err := tools.GetIntFromRequest(r, "id")
	if err != nil {
		return tools.JSONError("Invalid number for id")
	} else if id < 0 {
		errStr := fmt.Sprintf("Invalid number for id (%d)", id)
		return tools.JSONError(errStr)
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id,message,time,name,status FROM messages WHERE id = %d", id)
	rows, err := db.Query(sqlReq)
	if err != nil {
		return tools.JSONError(err.Error())
	}

	if rows.Next() {
		// Get infos
		var id int
		var txt, time, name, status string
		if err := rows.Scan(&id, &txt, &time, &name, &status); err != nil {
			return tools.JSONError(err.Error())
		}

		// Treat time to convenient format
		time, err := tools.SQLFormatDateTime(time)
		if err != nil {
			time = err.Error()
		}

		return messageToJSON(id, txt, time, name, status)
	}

	return tools.JSONError("No message corresponding to this ID")
}
