package api

import (
	"fmt"
	"math/rand"
	"net/http"
)

const authTokenSize = 20

type tokenBundle struct {
	bundle
}

func newTokenBundle(b bundle) tokenBundle {
	return tokenBundle{b}
}

func (t *tokenBundle) Process(r *http.Request) string {
	switch r.Method {
	case getMethod:
		return t.Get(r)
	case postMethod:
		return jsonError("Method not implemented")
	case putMethod:
		return jsonError("Method not implemented")
	case deleteMethod:
		return jsonError("Method not implemented")
	default:
		return jsonError("Unknown HTTP Method")
	}
}

func (t *tokenBundle) Get(r *http.Request) string {
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
			token, err := t.generatetokenBundle(id)
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

func (t *tokenBundle) generatetokenBundle(id int) (string, error) {
	// Generate the token
	b := make([]byte, authTokenSize)
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
