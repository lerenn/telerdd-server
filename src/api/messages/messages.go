package messages

import (
	"database/sql"
	"net/http"

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
	id       *ID
	send     *Send
	status   *Status
}

func New(data *data.Data, db *sql.DB, logger *log.Log) *Messages {
	var m Messages
	m.data = data
	m.db = db
	m.logger = logger
	m.next = NewNext(data, db, logger)
	m.previous = NewPrevious(data, db, logger)
	m.message = NewMessage(data, db, logger)
	m.id = NewID(data, db, logger)
	m.send = NewSend(data, db, logger)
	m.status = NewStatus(data, db, logger)
	return &m
}

func (m *Messages) Process(request string, r *http.Request) string {
	base, _ := tools.Split(request, "/")

	if base == "next" {
		return m.next.Process(r)
	} else if base == "previous" {
		return m.previous.Process(r)
	} else if base == "message" {
		return m.message.Process(r)
	} else if base == "id" {
		return m.id.Process(r)
	} else if base == "send" {
		return m.send.Process(r)
	} else if base == "status" {
		return m.status.Process(r)
	} else {
		return tools.JSONBadURL(r)
	}
}
