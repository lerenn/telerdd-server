package common

import (
	"net/http"
)

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
