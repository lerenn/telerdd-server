package messages

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type Message struct {
	// Infos
	data   *data.Data
	db     *sql.DB
	logger *log.Log
}

func NewMessage(data *data.Data, db *sql.DB, logger *log.Log) *Message {
	var m Message
	m.data = data
	m.db = db
	m.logger = logger
	return &m
}

func (m *Message) Process(r *http.Request) string {
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
	rows, err := m.db.Query(sqlReq)
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

	return tools.JSONError("No message corresponding to this Message")
}
