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
	message *Message
}

func New(data *data.Data, db *sql.DB, logger *log.Log) *Messages {
	var m Messages
	m.data = data
	m.db = db
	m.logger = logger
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
			return m.Post(r)
		case "PUT":
			return tools.JSONError("Method not implemented")
		case "DELETE":
			return tools.JSONError("Method not implemented")
		default:
			return tools.JSONError("Unknown HTTP Method")
		}
	} else if base == "message" {
		return m.message.Process(r)
	} else {
		return tools.JSONBadURL(r)
	}
}

func (m *Messages) Get(r *http.Request) string {
	// Get arguments
	status := getStatus(r)
	requestArgs := fmt.Sprintf("status REGEXP %q", status)

	start, present, err := tools.GetIntFromRequest(r, "start")
	if err != nil {
		return tools.JSONError("Error in 'start' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s AND id >= %d", requestArgs, start)
	}
	stop, present, err := tools.GetIntFromRequest(r, "stop")
	if err != nil {
		return tools.JSONError("Error in 'stop' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s AND id <= %d", requestArgs, stop)
	}
	offset, present, err := tools.GetIntFromRequest(r, "offset")
	if err != nil {
		return tools.JSONError("Error in 'offset' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s LIMIT %d", requestArgs, offset)
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id FROM messages WHERE %s", requestArgs)
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

	return "{\"messages\":[" + response + "]}"
}

func (m *Messages) Post(r *http.Request) string {
	// Check if authorized
	ip := tools.GetIp(r)
	t, err := m.data.ProceedMessageLimit(ip)
	if err != nil {
		return tools.JSONError("Error when check older messages: " + err.Error())
	} else if t != -1 {
		errStr := fmt.Sprintf("You already sent (or tried to send) a message %d seconds ago (from %s). Please wait.", t, ip)
		return tools.JSONError(errStr)
	}

	// Get message from request and format
	message := r.FormValue("message")
	if message == "" {
		return tools.JSONError("No text in your message")
	}
	message = tools.ReplaceBadCharacters(message)

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	}

	// Add to database
	stmt, errPrep := m.db.Prepare("INSERT messages SET ip=?,time=?,message=?,name=?,status=?")
	if errPrep != nil {
		return tools.JSONError(errPrep.Error())
	}

	_, errExec := stmt.Exec(ip, tools.SQLTimeNow(), message, name, "pending")
	if errExec != nil {
		return tools.JSONError(errExec.Error())
	}

	// Elaborate response
	m.logger.Print("Message posted (from " + ip + ") : " + message)
	return tools.JSONResponseOk()
}
