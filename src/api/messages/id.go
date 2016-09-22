package messages

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lerenn/telerdd-server/src/tools"
)

func id(db *sql.DB, r *http.Request) string {
	// Get status and start/stop
	status := getStatus(r)
	start, err := tools.GetIntFromRequest(r, "start")
	if err != nil || start < 0 {
		start = 0
	}
	stop, err := tools.GetIntFromRequest(r, "stop")
	if err != nil || stop < 0 {
		stop = 1000000
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id FROM messages WHERE status REGEXP %q AND (id BETWEEN %d AND %d)", status, start, stop)
	rows, err := db.Query(sqlReq)
	if err != nil {
		return tools.JSONError(err.Error())
	}

	var response string
	for i := 0; rows.Next(); i++ {
		// Get infos
		var id int
		if err := rows.Scan(&id); err != nil {
			return tools.JSONError(err.Error())
		}

		// Add to payload (older first)
		if i != 0 {
			response += ","
		}
		response += strconv.Itoa(id)
	}

	return "{\"id\":[" + response + "]}"
}
