package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lerenn/log"
)

type Messages struct {
	// Infos
	data   *Data
	db     *sql.DB
	logger *log.Log
	// API
	msg *Message
}

func newMessages(data *Data, db *sql.DB, logger *log.Log) *Messages {
	var m Messages
	m.data = data
	m.db = db
	m.logger = logger
	m.msg = newMessage(data, db, logger)
	return &m
}

func (m *Messages) Process(request string, r *http.Request) string {
	base, extend := splitString(request, "/")

	// If nothing more
	if base == "" {
		switch r.Method {
		case "GET":
			return m.Get(r)
		case "POST":
			return m.Post(r)
		case "PUT":
			return jsonError("Method not implemented")
		case "DELETE":
			return jsonError("Method not implemented")
		default:
			return jsonError("Unknown HTTP Method")
		}
	}

	// Let message take care of the request : get id
	id, err := strconv.Atoi(base)
	if err != nil || id < 0 {
		errStr := fmt.Sprintf("Invalid message id : %q", base)
		return jsonError(errStr)
	}
	return m.msg.Process(r, id, extend)
}

func (m *Messages) Get(r *http.Request) string {
	// Get arguments
	status := getStatus(r)
	requestArgs := fmt.Sprintf("status REGEXP %q", status)

	start, present, err := getRequestInt(r, "start")
	if err != nil {
		return jsonError("Error in 'start' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s AND id >= %d", requestArgs, start)
	}
	stop, present, err := getRequestInt(r, "stop")
	if err != nil {
		return jsonError("Error in 'stop' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s AND id <= %d", requestArgs, stop)
	}
	offset, present, err := getRequestInt(r, "offset")
	if err != nil {
		return jsonError("Error in 'offset' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s LIMIT %d", requestArgs, offset)
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id FROM messages WHERE %s", requestArgs)
	rows, err := m.db.Query(sqlReq)
	if err != nil {
		return jsonError(err.Error())
	}
	defer rows.Close()

	var response string
	for i := 0; rows.Next(); i++ {
		// Get infos
		var list int
		if err := rows.Scan(&list); err != nil {
			return jsonError(err.Error())
		}

		// Add to payload (older first)
		if i != 0 {
			response += ","
		}
		response += strconv.Itoa(list)
	}

	return "{\"messages\":[" + response + "]}"
}

func (m *Messages) Post(r *http.Request) string {
	// Check if authorized
	ip := getRequestIP(r)
	t, err := m.data.ProceedMessageLimit(ip)
	if err != nil {
		return jsonError("Error when check older messages: " + err.Error())
	} else if t != -1 {
		errStr := fmt.Sprintf("You already sent (or tried to send) a message %d seconds ago (from %s). Please wait.", t, ip)
		return jsonError(errStr)
	}

	// Get message from request and format
	message := r.FormValue("message")
	if message == "" {
		return jsonError("No text in your message")
	}
	message = replaceBadChar(message)

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	}

	// Add to database
	stmt, err := m.db.Prepare("INSERT messages SET ip=?,time=?,message=?,img=?,name=?,status=?")
	if err != nil {
		return jsonError(err.Error())
	}

	res, err := stmt.Exec(ip, sqlTimeNow(), message, false, name, "pending")
	if err != nil {
		return jsonError(err.Error())
	}

	// Get id
	id, err := res.LastInsertId()
	if err != nil {
		return jsonError(err.Error())
	}

	// Elaborate response
	m.logger.Print("Message posted (from " + ip + ") : " + message)
	return fmt.Sprintf("{\"id\":%d}", id)
}

// Private functions
////////////////////////////////////////////////////////////////////////////////

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
