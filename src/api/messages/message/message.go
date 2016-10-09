package message

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
)

type Message struct {
	data   *common.Data
	db     *sql.DB
	logger *log.Log
	// API
	img *Image
}

func New(data *common.Data, db *sql.DB, logger *log.Log) *Message {
	var m Message
	m.data = data
	m.db = db
	m.logger = logger
	m.img = NewImage(data, db, logger)
	return &m
}

func (m *Message) Process(r *http.Request, id int, request string) string {
	base, _ := common.Split(request, "/")

	// If there is nothing
	if base == "" {
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
	} else if base == "image" {
		return m.img.Process(r, id)
	} else {
		return common.JSONBadURL(r)
	}
}

func (m *Message) Get(r *http.Request, id int) string {
	sqlReq := fmt.Sprintf("SELECT id,time,message,img,name,status FROM messages WHERE id = %d", id)
	rows, err := m.db.Query(sqlReq)
	if err != nil {
		return common.JSONError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get infos
		var id int
		var img bool
		var txt, time, name, status string
		if err := rows.Scan(&id, &time, &txt, &img, &name, &status); err != nil {
			return common.JSONError(err.Error())
		}

		// Treat time to convenient format
		time, err := common.SQLFormatDateTime(time)
		if err != nil {
			time = err.Error()
		}

		return common.MessageToJSON(id, img, txt, time, name, status)
	}

	return common.JSONError("No message corresponding to this ID")
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

	// Exec request
	_, errExec := stmt.Exec(status, id)
	if errExec != nil {
		return common.JSONError(errExec.Error())
	}

	return common.JSONResponseOk()
}
