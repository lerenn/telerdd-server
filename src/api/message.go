package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lerenn/log"
)

type Message struct {
	data   *Data
	db     *sql.DB
	logger *log.Log
	// API
	img *Image
}

func newMessage(data *Data, db *sql.DB, logger *log.Log) *Message {
	var m Message
	m.data = data
	m.db = db
	m.logger = logger
	m.img = newImage(data, db, logger)
	return &m
}

func (m *Message) Process(r *http.Request, id int, request string) string {
	base, _ := splitString(request, "/")

	// If there is nothing
	if base == "" {
		switch r.Method {
		case "GET":
			return m.Get(r, id)
		case "POST":
			return jsonError("Method not implemented")
		case "PUT":
			return m.Put(r, id)
		case "DELETE":
			return jsonError("Method not implemented")
		default:
			return jsonError("Unknown HTTP Method")
		}
	} else if base == "image" {
		return m.img.Process(r, id)
	} else {
		return jsonBadURL(r)
	}
}

func (m *Message) Get(r *http.Request, id int) string {
	sqlReq := fmt.Sprintf("SELECT id,time,message,img,name,status FROM messages WHERE id = %d", id)
	rows, err := m.db.Query(sqlReq)
	if err != nil {
		return jsonError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get infos
		var id int
		var img bool
		var txt, time, name, status string
		if err := rows.Scan(&id, &time, &txt, &img, &name, &status); err != nil {
			return jsonError(err.Error())
		}

		// Treat time to convenient format
		time, err := sqlFormatDateTime(time)
		if err != nil {
			time = err.Error()
		}

		return MessageToJSON(id, img, txt, time, name, status)
	}

	return jsonError("No message corresponding to this ID")
}

func (m *Message) Put(r *http.Request, id int) string {
	// Check permission
	errAuth := authorized(m.db, r, 1)
	if errAuth != nil {
		return jsonError(errAuth.Error())
	}

	// Get status
	status := r.FormValue("status")
	if status != "accepted" && status != "refused" {
		return jsonError("Invalid status")
	}

	// Prepare request
	stmt, errPrep := m.db.Prepare("UPDATE messages SET status=? WHERE id=?")
	if errPrep != nil {
		return jsonError(errPrep.Error())
	}

	// Exec request
	_, errExec := stmt.Exec(status, id)
	if errExec != nil {
		return jsonError(errExec.Error())
	}

	return jsonResponseOk()
}

func MessageToJSON(id int, img bool, txt, time, name, status string) string {
	// Check image bool
	imgTxt := "false"
	if img {
		imgTxt = "true"
	}

	// Form message
	response := "{"
	response += "\"id\":\"" + strconv.Itoa(id) + "\","
	response += "\"text\":\"" + txt + "\","
	response += "\"img\":\"" + imgTxt + "\","
	response += "\"time\":\"" + time + "\","
	response += "\"name\":\"" + name + "\","
	response += "\"status\":\"" + status + "\""
	response += "}"
	return response
}
