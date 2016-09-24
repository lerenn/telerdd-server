package messages

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type Messages struct {
	// Infos
	data   *data.Data
	db     *sql.DB
	logger *log.Log
	// API
	next     *Next
	previous *Previous
	message  *Message
}

func New(data *data.Data, db *sql.DB, logger *log.Log) *Messages {
	var m Messages
	m.data = data
	m.db = db
	m.logger = logger
	m.next = NewNext(data, db, logger)
	m.previous = NewPrevious(data, db, logger)
	m.message = NewMessage(data, db, logger)
	return &m
}

func (m *Messages) Process(request string, r *http.Request) string {
	base, _ := tools.Split(request, "/")

	if base == "" {
		switch r.Method {
		case "GET":
			return m.Get(r)
		case "POST":
			return tools.JSONError("Method not implemented")
		case "PUT":
			return tools.JSONError("Method not implemented")
		case "DELETE":
			return tools.JSONError("Method not implemented")
		default:
			return tools.JSONError("Unknown HTTP Method")
		}
	} else if base == "next" {
		return m.next.Process(r)
	} else if base == "previous" {
		return m.previous.Process(r)
	} else if base == "message" {
		return m.message.Process(r)
	} else {
		return tools.JSONBadURL(r)
	}
}

func (m *Messages) Get(r *http.Request) string {
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
	rows, err := m.db.Query(sqlReq)
	if err != nil {
		return tools.JSONError(err.Error())
	}

	var response string
	for i := 0; rows.Next(); i++ {
		// Get infos
		var list int
		if err := rows.Scan(&list); err != nil {
			return tools.JSONError(err.Error())
		}

		// Add to payload (older first)
		if i != 0 {
			response += ","
		}
		response += strconv.Itoa(list)
	}

	return "{\"list\":[" + response + "]}"
}
