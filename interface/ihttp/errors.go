package ihttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/utrack/woofer/lib/bizerr"
)

type httpError struct {
	Error string
	Stack string
}

func renderError(w http.ResponseWriter, err error, retCode int) {
	if err == nil {
		return
	}
	if code := bizerr.Type(err); code != bizerr.ErrorUnknown {
		switch code {
		case bizerr.ErrorUserInput:
			retCode = 400
		case bizerr.ErrorConflict:
			retCode = http.StatusConflict
		}
	}

	w.WriteHeader(retCode)
	e := httpError{
		Error: err.Error(),
		Stack: fmt.Sprintf("%+v", err),
	}
	json.NewEncoder(w).Encode(e)
}
