package messages

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type Previous struct {
	// Infos
	data   *data.Data
	db     *sql.DB
	logger *log.Log
}

func NewPrevious(data *data.Data, db *sql.DB, logger *log.Log) *Previous {
	var p Previous
	p.data = data
	p.db = db
	p.logger = logger
	return &p
}

func (p *Previous) Process(r *http.Request) string {
	// Get status
	status := getStatus(r)

	// Get last id
	id, err := tools.GetIntFromRequest(r, "id")
	if err != nil {
		return tools.JSONError("Invalid number for id")
	} else if id < 0 {
		id = 0
	}

	// Get offset
	offset, err := tools.GetIntFromRequest(r, "offset")
	if err != nil {
		return tools.JSONError("Invalid number for offset")
	} else if offset < 0 {
		errStr := fmt.Sprintf("Invalid number for offset (%d)", offset)
		return tools.JSONError(errStr)
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT id,message,time,name,status FROM messages WHERE status REGEXP %q AND id < %d ORDER BY id DESC LIMIT %d", status, id, offset)
	rows, err := p.db.Query(sqlReq)
	if err != nil {
		return tools.JSONError(err.Error())
	}

	return createListFromSQLRequest(rows)
}
