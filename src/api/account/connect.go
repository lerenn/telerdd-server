package account

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd/src/data"
	"github.com/lerenn/telerdd/src/tools"
)

func connect(infos *data.Data, db *sql.DB, l *log.Log, r *http.Request) string {
	// Get infos
	username := r.FormValue("username")
	psswd := r.FormValue("password")
	accountType, errType := tools.GetIntFromRequest(r, "type")
	if errType != nil {
		return tools.JSONError(errType.Error())
	}

	// Get user infos from db
	sqlReq := fmt.Sprintf("SELECT id, password, type FROM users WHERE username = %q", username)
	rows, err := db.Query(sqlReq)
	if err != nil {
		return tools.JSONError(err.Error())
	}

	if rows.Next() {
		// Get password
		var id, accountTypeFromDB int
		var psswdFromDB string
		if err := rows.Scan(&id, &psswdFromDB, &accountTypeFromDB); err != nil {
			return tools.JSONError(err.Error())
		}

		if psswd == psswdFromDB && accountType >= accountTypeFromDB {
			token, err := generateToken(db, id)
			if err != nil {
				tools.JSONError("Error when generating token : " + err.Error())
			}

			return fmt.Sprintf("{\"token\": %q}", token)
		} else if accountType < accountTypeFromDB {
			return tools.JSONError("Not authorized")
		}
	}

	return tools.JSONError("Account or password incorrect")
}

func generateToken(db *sql.DB, id int) (string, error) {
	// Generate the token
	b := make([]byte, TOKEN_SIZE)
	rand.Read(b)
	token := fmt.Sprintf("%x", b)

	// Add it to the user in db
	stmt, errPrep := db.Prepare("UPDATE users SET token=? WHERE id=?")
	if errPrep != nil {
		return "", errPrep
	}

	_, errExec := stmt.Exec(token, id)
	if errExec != nil {
		return "", errExec
	}

	return token, nil
}
