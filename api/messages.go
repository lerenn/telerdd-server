package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type messagesBundle struct {
	bundle
	msg messageBundle
}

func newMessagesBundle(b bundle) messagesBundle {
	return messagesBundle{b, newMessageBundle(b)}
}

func (m *messagesBundle) Process(request string, r *http.Request) string {
	base, extend := splitString(request, "/")

	// If nothing more
	if base == "" {
		switch r.Method {
		case getMethod:
			return m.Get(r)
		case postMethod:
			return m.Post(r)
		case putMethod:
			return jsonError("Method not implemented")
		case deleteMethod:
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

func (m *messagesBundle) Get(r *http.Request) string {
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

func (m *messagesBundle) Post(r *http.Request) string {
	var status string

	// Check if authorized
	ip := getRequestIP(r)
	t, err := m.data.ProceedMessageLimit(ip)
	if err != nil {
		return jsonError("Error when check older messages: " + err.Error())
	} else if t != -1 {
		errStr := fmt.Sprintf("You already sent (or tried to send) a message %d seconds ago (from %s). Please wait.", t, ip)
		return jsonError(errStr)
	}

	// Get name from request
	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	} else {
		name = template.HTMLEscapeString(name)
	}

	// Get img from request
	img := r.FormValue("image")
	imgPresence := strings.Contains(img, "base64") || strings.Contains(img, "http")
	// Process img
	if imgPresence {
		if img, err = processImg(m.data, img); err != nil {
			return jsonError(err.Error())
		}
	}

	// Get message from request
	message := r.FormValue("message")
	if message == "" && imgPresence == false {
		return jsonError("No text in your message")
	} else if msgLen := len(message); msgLen > m.data.MsgLimitSize {
		errTxt := fmt.Sprintf("Message limit size exceeded (%d characters for %d characters max)", msgLen, m.data.MsgLimitSize)
		return jsonError(errTxt)
	} else if message != "" {
		message = template.HTMLEscapeString(message)
		message = replaceBadChar(message)
	}

	// Add to database
	stmt, err := m.db.Prepare("INSERT messages SET ip=?,time=?,message=?,img=?,name=?,status=?")
	if err != nil {
		return jsonError(err.Error())
	}

	if (imgPresence && m.data.MsgModerationWithImg) || (!imgPresence && m.data.MsgModerationWithoutImg) {
		status = "pending"
	} else {
		status = "accepted"
	}

	res, err := stmt.Exec(ip, sqlTimeNow(), message, imgPresence, name, status)
	if err != nil {
		return jsonError(err.Error())
	}

	// Get id
	id, err := res.LastInsertId()
	if err != nil {
		return jsonError(err.Error())
	}

	// Save image if there is one
	if imgPresence {
		if err := saveImg(m.db, img, int(id)); err != nil {
			return jsonError(err.Error())
		}
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
