package messages

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
)

type Message struct {
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

func (m *Message) Process(r *http.Request, request string) string {
	base, _ := common.Split(request, "/")

	// Convert to number
	id, err := strconv.Atoi(base)
	if err != nil || id < 0 {
		errStr := fmt.Sprintf("Invalid message id : %q", base)
		return common.JSONError(errStr)
	}

	switch r.Method {
	case "GET":
		return m.Get(r, id)
	case "POST":
		return common.JSONError("Method not implemented")
	case "PUT":
		return m.Put(r, id)
	case "DELETE":
		return common.JSONError("Method not implemented")
	default:
		return common.JSONError("Unknown HTTP Method")
	}
}

func (m *Message) Get(r *http.Request, id int) string {
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

func (m *Message) Put(r *http.Request, id int) string {
	// Check permission
	errAuth := common.Authorized(m.db, r, 1)
	if errAuth != nil {
		return common.JSONError(errAuth.Error())
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
