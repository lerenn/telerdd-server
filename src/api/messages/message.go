package messages

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
)

type Message struct {
	// Infos
	data   *common.Data
	db     *sql.DB
	logger *log.Log
}

func NewMessage(data *common.Data, db *sql.DB, logger *log.Log) *Message {
	var m Message
	m.data = data
	m.db = db
	m.logger = logger
	return &m
}

func (m *Message) Process(r *http.Request) string {
	switch r.Method {
	case "GET":
		return m.Get(r)
	case "POST":
		return common.JSONError("Method not implemented")
	case "PUT":
		return m.Put(r)
	case "DELETE":
		return common.JSONError("Method not implemented")
	default:
		return common.JSONError("Unknown HTTP Method")
	}
}

func (m *Message) Get(r *http.Request) string {
	// Get id
	id, _, err := common.GetIntFromRequest(r, "id")
	if err != nil {
		return common.JSONError("Invalid number for id")
	} else if id < 0 {
		errStr := fmt.Sprintf("Invalid number for id (%d)", id)
		return common.JSONError(errStr)
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id,message,time,name,status FROM messages WHERE id = %d", id)
	rows, err := m.db.Query(sqlReq)
	if err != nil {
		return common.JSONError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get infos
		var id int
		var txt, time, name, status string
		if err := rows.Scan(&id, &txt, &time, &name, &status); err != nil {
			return common.JSONError(err.Error())
		}

		// Treat time to convenient format
		time, err := common.SQLFormatDateTime(time)
		if err != nil {
			time = err.Error()
		}

		return messageToJSON(id, txt, time, name, status)
	}

	return common.JSONError("No message corresponding to this Message")
}

func (m *Message) Put(r *http.Request) string {
	// Check permission
	errAuth := common.Authorized(m.db, r, 1)
	if errAuth != nil {
		return common.JSONError(errAuth.Error())
	}

	// Get id
	id, _, err := common.GetIntFromRequest(r, "id")
	if err != nil {
		return common.JSONError("Invalid number for id")
	} else if id < 0 {
		id = 0
	}

	// Get status
	status := r.FormValue("status")
	if status != "accepted" && status != "refused" {
		return common.JSONError("Invalid status")
	}

	// Prepare request
	stmt, errPrep := m.db.Prepare("UPDATE messages SET status=? WHERE id=?")
	if errPrep != nil {
		return common.JSONError(errPrep.Error())
	}

	_, errExec := stmt.Exec(status, id)
	if errExec != nil {
		return common.JSONError(errExec.Error())
	}

	return common.JSONResponseOk()
}
