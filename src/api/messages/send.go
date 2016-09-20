package messages

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lerenn/log"
	"github.com/lerenn/telerdd/src/data"
	"github.com/lerenn/telerdd/src/tools"
)

func send(infos *data.Data, db *sql.DB, l *log.Log, r *http.Request) string {
	// Check if authorized
	ip := tools.GetIp(r)
	t, err := infos.ProceedMessageLimit(ip)
	if err != nil {
		return tools.JSONError("Error when check older messages: " + err.Error())
	} else if t != -1 {
		errStr := fmt.Sprintf("You already sent (or tried to send) a message %d seconds ago (from %s). Please wait.", t, ip)
		return tools.JSONError(errStr)
	}

	// Get message from request and format
	message := r.FormValue("message")
	if message == "" {
		return tools.JSONError("No text in your message")
	}
	message = tools.ReplaceBadCharacters(message)

	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	}

	// Add to database
	stmt, errPrep := db.Prepare("INSERT messages SET ip=?,time=?,message=?,name=?,status=?")
	if errPrep != nil {
		return tools.JSONError(errPrep.Error())
	}

	_, errExec := stmt.Exec(ip, tools.SQLTimeNow(), message, name, "pending")
	if errExec != nil {
		return tools.JSONError(errExec.Error())
	}

	// Elaborate response
	l.Print("Message posted (from " + ip + ") : " + message)
	return tools.JSONResponseOk()
}
