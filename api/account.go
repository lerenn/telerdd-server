package api

import "net/http"

type accountBundle struct {
	bundle
	token tokenBundle
}

func newAccountBundle(b bundle) accountBundle {
	return accountBundle{b, newTokenBundle(b)}
}

func (a *accountBundle) Process(request string, r *http.Request) string {
	base, _ := splitString(request, "/")

	if base == "token" {
		return a.token.Process(r)
	}
	return jsonBadURL(r)
}
