package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func authorized(db *sql.DB, r *http.Request, accountType int) error {
	// Get username
	username := r.FormValue("username")
	if username == "" {
		return errors.New("No username provided for auth")
	}

	// Get token
	token := r.FormValue("token")
	if token == "" {
		return errors.New("No token provided for auth")
	}

	// Get complete list
	sqlReq := fmt.Sprintf("SELECT token, type FROM users WHERE username = %q", username)
	rows, err := db.Query(sqlReq)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var accountTypeFromDB int
		var tokenFromDB string
		if err := rows.Scan(&tokenFromDB, &accountTypeFromDB); err != nil {
			return err
		}

		if token == tokenFromDB && accountType >= accountTypeFromDB {
			return nil
		} else if accountType < accountTypeFromDB {
			return errors.New("Not authorized")
		}
		// test := fmt.Sprintf("%q, %q", token, tokenFromDB)
		return errors.New("Invalid token")
	}

	return errors.New("Invalid username")
}
