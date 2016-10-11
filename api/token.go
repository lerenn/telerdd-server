package api

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/lerenn/log"
	cst "github.com/lerenn/telerdd-server/constants"
)

type Token struct {
	data   *Data
	db     *sql.DB
	logger *log.Log
}

func newToken(data *Data, db *sql.DB, logger *log.Log) *Token {
	var t Token
	t.data = data
	t.db = db
	t.logger = logger
	return &t
}

func (t *Token) Process(r *http.Request) string {
	switch r.Method {
	case "GET":
		return t.Get(r)
	case "POST":
		return jsonError("Method not implemented")
	case "PUT":
		return jsonError("Method not implemented")
	case "DELETE":
		return jsonError("Method not implemented")
	default:
		return jsonError("Unknown HTTP Method")
	}
}

func (t *Token) Get(r *http.Request) string {
	// Get infos
	username := r.FormValue("username")
	psswd := r.FormValue("password")
	accountType, _, errType := getRequestInt(r, "type")
	if errType != nil {
		return jsonError(errType.Error())
	}

	// Get user infos from db
	sqlReq := fmt.Sprintf("SELECT id, password, type FROM users WHERE username = %q", username)
	rows, err := t.db.Query(sqlReq)
	if err != nil {
		return jsonError(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		// Get password
		var id, accountTypeFromDB int
		var psswdFromDB string
		if err := rows.Scan(&id, &psswdFromDB, &accountTypeFromDB); err != nil {
			return jsonError(err.Error())
		}

		if psswd == psswdFromDB && accountType >= accountTypeFromDB {
			token, err := t.generateToken(id)
			if err != nil {
				jsonError("Error when generating token : " + err.Error())
			}

			return fmt.Sprintf("{\"token\": %q}", token)
		} else if accountType < accountTypeFromDB {
			return jsonError("Not authorized")
		}
	}

	return jsonError("Account or password incorrect")
}

// Private methods
////////////////////////////////////////////////////////////////////////////////

func (t *Token) generateToken(id int) (string, error) {
	// Generate the token
	b := make([]byte, cst.AUTH_TOKEN_SIZE)
	rand.Read(b)
	token := fmt.Sprintf("%x", b)

	// Add it to the user in db
	stmt, errPrep := t.db.Prepare("UPDATE users SET token=? WHERE id=?")
	if errPrep != nil {
		return "", errPrep
	}

	_, errExec := stmt.Exec(token, id)
	if errExec != nil {
		return "", errExec
	}

	return token, nil
}
