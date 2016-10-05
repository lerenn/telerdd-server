package message

import (
	"database/sql"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
)

type Image struct {
	data   *common.Data
	db     *sql.DB
	logger *log.Log
}

func NewImage(data *common.Data, db *sql.DB, logger *log.Log) *Image {
	var i Image
	i.data = data
	i.db = db
	i.logger = logger
	return &i
}

func (i *Image) Process(r *http.Request, id int) string {
	switch r.Method {
	case "GET":
		return common.JSONError("Method not implemented")
	case "POST":
		return i.Post(r, id)
	case "PUT":
		return common.JSONError("Method not implemented")
	case "DELETE":
		return common.JSONError("Method not implemented")
	default:
		return common.JSONError("Unknown HTTP Method")
	}
}

func (i *Image) Post(r *http.Request, id int) string {
	// Get image from request
	image := r.FormValue("image")
	if image == "" {
		return common.JSONError("No image provided")
	}

	// TODO: Check if there is no image

	// Prepare add to database
	stmt, err := i.db.Prepare("INSERT images SET time=?,msg_id=?,img=?")
	if err != nil {
		return common.JSONError(err.Error())
	}

	// Exec request
	_, err = stmt.Exec(common.SQLTimeNow(), id, image)
	if err != nil {
		return common.JSONError(err.Error())
	}

	// Prepare request
	stmt, errPrep := i.db.Prepare("UPDATE messages SET img=? WHERE id=?")
	if errPrep != nil {
		return common.JSONError(errPrep.Error())
	}

	// Exec request
	_, errExec := stmt.Exec(true, id)
	if errExec != nil {
		return common.JSONError(errExec.Error())
	}

	return common.JSONResponseOk()
}
