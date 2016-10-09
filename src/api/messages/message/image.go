package message

import (
	"database/sql"
	"fmt"
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
		return i.Get(r, id)
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

func (i *Image) Get(r *http.Request, id int) string {
	sqlReq := fmt.Sprintf("SELECT img FROM images WHERE msg_id = %d", id)
	rows, err := i.db.Query(sqlReq)
	if err != nil {
		return common.JSONError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get img
		var img string
		if err := rows.Scan(&img); err != nil {
			return common.JSONError(err.Error())
		}

		return fmt.Sprintf("{\"img\":%q}", img)
	}

	return common.JSONError("No image corresponding to this message ID")
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
