package messages

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
)

type Messages struct {
	// Infos
	data   *common.Data
	db     *sql.DB
	logger *log.Log
	// API
	message *Message
}

func New(data *common.Data, db *sql.DB, logger *log.Log) *Messages {
	var m Messages
	m.data = data
	m.db = db
	m.logger = logger
	m.message = NewMessage(data, db, logger)
	return &m
}

func (m *Messages) Process(request string, r *http.Request) string {
	base, _ := common.Split(request, "/")

	if base == "" {
		switch r.Method {
		case "GET":
			return m.Get(r)
		case "POST":
			return m.Post(r)
		case "PUT":
			return common.JSONError("Method not implemented")
		case "DELETE":
			return common.JSONError("Method not implemented")
		default:
			return common.JSONError("Unknown HTTP Method")
		}
	} else {
		return m.message.Process(r, request)
	}
}

func (m *Messages) Get(r *http.Request) string {
	// Get arguments
	status := getStatus(r)
	requestArgs := fmt.Sprintf("status REGEXP %q", status)

	start, present, err := common.GetIntFromRequest(r, "start")
	if err != nil {
		return common.JSONError("Error in 'start' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s AND id >= %d", requestArgs, start)
	}
	stop, present, err := common.GetIntFromRequest(r, "stop")
	if err != nil {
		return common.JSONError("Error in 'stop' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s AND id <= %d", requestArgs, stop)
	}
	offset, present, err := common.GetIntFromRequest(r, "offset")
	if err != nil {
		return common.JSONError("Error in 'offset' argument")
	} else if present {
		requestArgs = fmt.Sprintf("%s LIMIT %d", requestArgs, offset)
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id FROM messages WHERE %s", requestArgs)
	rows, err := m.db.Query(sqlReq)
	if err != nil {
		return common.JSONError(err.Error())
	}
	defer rows.Close()

	var response string
	for i := 0; rows.Next(); i++ {
		// Get infos
		var list int
		if err := rows.Scan(&list); err != nil {
			return common.JSONError(err.Error())
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
	ip := common.GetIp(r)
	t, err := m.data.ProceedMessageLimit(ip)
	if err != nil {
		return common.JSONError("Error when check older messages: " + err.Error())
	} else if t != -1 {
		errStr := fmt.Sprintf("You already sent (or tried to send) a message %d seconds ago (from %s). Please wait.", t, ip)
		return common.JSONError(errStr)
	}

	// Get message from request and format
	message := r.FormValue("message")
	if message == "" {
		return common.JSONError("No text in your message")
	}
	message = common.ReplaceBadCharacters(message)

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	}

	// Add to database
	stmt, errPrep := m.db.Prepare("INSERT messages SET ip=?,time=?,message=?,name=?,status=?")
	if errPrep != nil {
		return common.JSONError(errPrep.Error())
	}

	_, errExec := stmt.Exec(ip, common.SQLTimeNow(), message, name, "pending")
	if errExec != nil {
		return common.JSONError(errExec.Error())
	}

	// Elaborate response
	m.logger.Print("Message posted (from " + ip + ") : " + message)
	return common.JSONResponseOk()
}
