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

func GetIntFromRequest(r *http.Request, name string) (int, bool, error) {
	// Get string
	nbrStr := r.FormValue(name)

	// If string is empty, then it is not passed in parameter
	if nbrStr == "" {
		return 0, false, nil
	}

	// Try to change string to nbr
	nbr, err := strconv.Atoi(nbrStr)
	if err != nil {
		errStr := fmt.Sprintf("Invalid number : %q", nbrStr)
		return 0, true, errors.New(errStr)
	}

	// Return number
	return nbr, true, nil
}
