package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrInvalidCpntentType = errors.New("invalid content type")
)

func RequestJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrInvalidCpntentType
	}
	return json.NewDecoder(r.Body).Decode(dst)
}
