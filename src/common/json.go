package common

import (
	"net/http"
	"strconv"
)

func MessageToJSON(id int, img bool, txt, time, name, status string) string {
	// Check image bool
	imgTxt := "false"
	if img {
		imgTxt = "true"
	}

	// Form message
	response := "{"
	response += "\"id\":\"" + strconv.Itoa(id) + "\","
	response += "\"text\":\"" + txt + "\","
	response += "\"img\":\"" + imgTxt + "\","
	response += "\"time\":\"" + time + "\","
	response += "\"name\":\"" + name + "\","
	response += "\"status\":\"" + status + "\""
	response += "}"
	return response
}

func JSONResponseOk() string {
	return "{\"response\": \"OK\"}"
}

func JSONError(txt string) string {
	return "{\"error\":\"" + ReplaceBadCharacters(txt) + "\"}"
}

func JSONBadURL(r *http.Request) string {
	url := r.URL.Path[1:]
	return JSONError("Bad URL: " + url)
}
