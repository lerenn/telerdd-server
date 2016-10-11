package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
)

type Image struct {
	data   *Data
	db     *sql.DB
	logger *log.Log
}

func newImage(data *Data, db *sql.DB, logger *log.Log) *Image {
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
		return jsonError("Method not implemented")
	case "DELETE":
		return jsonError("Method not implemented")
	default:
		return jsonError("Unknown HTTP Method")
	}
}

func (i *Image) Get(r *http.Request, id int) string {
	sqlReq := fmt.Sprintf("SELECT img FROM images WHERE msg_id = %d", id)
	rows, err := i.db.Query(sqlReq)
	if err != nil {
		return jsonError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get img
		var img string
		if err := rows.Scan(&img); err != nil {
			return jsonError(err.Error())
		}

		return fmt.Sprintf("{\"img\":%q}", img)
	}

	return jsonError("No image corresponding to this message ID")
}

func (i *Image) Post(r *http.Request, id int) string {
	// Get image from request
	image := r.FormValue("image")
	if image == "" {
		return jsonError("No image provided")
	}

	if err := saveImage(i.db,image,id); err != nil {
		return jsonError(err.Error())
	}
	return jsonResponseOk()
}

// Save image
func saveImage(db *sql.DB, img string, id int) error {
	// TODO: Check if there is already an image

	// Prepare add to database
	stmt, err := db.Prepare("INSERT images SET time=?,msg_id=?,img=?")
	if err != nil {
		return err
	}

	// Exec request
	_, err = stmt.Exec(sqlTimeNow(), id, img)
	if err != nil {
		return err
	}

	// Prepare request
	stmt, errPrep := db.Prepare("UPDATE messages SET img=? WHERE id=?")
	if errPrep != nil {
		return errPrep
	}

	// Exec request
	_, errExec := stmt.Exec(true, id)
	if errExec != nil {
		return errExec
	}

	return nil
}
