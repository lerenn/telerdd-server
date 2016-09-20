package tools

import (
	"errors"
	"fmt"
	"github.com/tomasen/realip"
	"net/http"
	"strconv"
)

func GetIp(r *http.Request) string {
	return realip.RealIP(r)
}

func GetIntFromRequest(r *http.Request, name string) (int, error) {
	nbrStr := r.FormValue(name)
	nbr, err := strconv.Atoi(nbrStr)
	if err != nil {
		errStr := fmt.Sprintf("Invalid number : %q", nbrStr)
		return 0, errors.New(errStr)
	}
	return nbr, nil
}
