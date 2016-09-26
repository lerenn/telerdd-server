package account

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/common"
	cst "github.com/lerenn/telerdd-server/src/constants"
)

type Token struct {
	data   *common.Data
	db     *sql.DB
	logger *log.Log
}

func NewToken(data *common.Data, db *sql.DB, logger *log.Log) *Token {
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
		return common.JSONError("Method not implemented")
	case "PUT":
		return common.JSONError("Method not implemented")
	case "DELETE":
		return common.JSONError("Method not implemented")
	default:
		return common.JSONError("Unknown HTTP Method")
	}
}

func (t *Token) Get(r *http.Request) string {
	// Get infos
	username := r.FormValue("username")
	psswd := r.FormValue("password")
	accountType, _, errType := common.GetIntFromRequest(r, "type")
	if errType != nil {
		return common.JSONError(errType.Error())
	}

	// Get user infos from db
	sqlReq := fmt.Sprintf("SELECT id, password, type FROM users WHERE username = %q", username)
	rows, err := t.db.Query(sqlReq)
	if err != nil {
		return common.JSONError(err.Error())
	}

	if rows.Next() {
		// Get password
		var id, accountTypeFromDB int
		var psswdFromDB string
		if err := rows.Scan(&id, &psswdFromDB, &accountTypeFromDB); err != nil {
			return common.JSONError(err.Error())
		}

		if psswd == psswdFromDB && accountType >= accountTypeFromDB {
			token, err := t.generateToken(id)
			if err != nil {
				common.JSONError("Error when generating token : " + err.Error())
			}

			return fmt.Sprintf("{\"token\": %q}", token)
		} else if accountType < accountTypeFromDB {
			return common.JSONError("Not authorized")
		}
	}

	return common.JSONError("Account or password incorrect")
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
