package messages

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/lerenn/telerdd/src/tools"
)

func messageToJSON(id int, txt, time, name, status string) string {
	response := "{"
	response += "\"id\":\"" + strconv.Itoa(id) + "\","
	response += "\"text\":\"" + txt + "\","
	response += "\"time\":\"" + time + "\","
	response += "\"name\":\"" + name + "\","
	response += "\"status\":\"" + status + "\""
	response += "}"
	return response
}

func createListFromSQLRequest(rows *sql.Rows) string {
	var response string
	for i := 0; rows.Next(); i++ {
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

		// Add to payload (older first)
		if i == 0 {
			response = messageToJSON(id, txt, time, name, status)
		} else {
			response += "," + messageToJSON(id, txt, time, name, status)
		}
	}

	return "{\"messages\":[" + response + "]}"
}

func getStatus(r *http.Request) string {
	status := r.FormValue("status")
	if status == "rejected" {
		return "rejected"
	} else if status == "pending" {
		return "pending"
	} else if status == "all" {
		return ".*"
	}
	return "accepted"
}
