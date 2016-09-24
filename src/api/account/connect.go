package account

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd-server/src/data"
	"github.com/lerenn/telerdd-server/src/tools"
)

type Connect struct {
	data   *data.Data
	db     *sql.DB
	logger *log.Log
}

func NewConnect(data *data.Data, db *sql.DB, logger *log.Log) *Connect {
	var c Connect
	c.data = data
	c.db = db
	c.logger = logger
	return &c
}

func (c *Connect) Process(r *http.Request) string {
	// Get infos
	username := r.FormValue("username")
	psswd := r.FormValue("password")
	accountType, errType := tools.GetIntFromRequest(r, "type")
	if errType != nil {
		return tools.JSONError(errType.Error())
	}

	// Get user infos from db
	sqlReq := fmt.Sprintf("SELECT id, password, type FROM users WHERE username = %q", username)
	rows, err := c.db.Query(sqlReq)
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
			token, err := c.generateToken(id)
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

func (c *Connect) generateToken(id int) (string, error) {
	// Generate the token
	b := make([]byte, TOKEN_SIZE)
	rand.Read(b)
	token := fmt.Sprintf("%x", b)

	// Add it to the user in db
	stmt, errPrep := c.db.Prepare("UPDATE users SET token=? WHERE id=?")
	if errPrep != nil {
		return "", errPrep
	}

	_, errExec := stmt.Exec(token, id)
	if errExec != nil {
		return "", errExec
	}

	return token, nil
}
