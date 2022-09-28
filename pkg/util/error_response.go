package util

import (
	"net/http"

	"fakedating/pkg/payload"
)

func WriteErrorResponse(msg string, err error, statusCode int, w http.ResponseWriter) {
	e := payload.ErrorResponse{Msg: msg}
	if err != nil {
		e.Error = err.Error()
	}

	WriteJSONResponse(e, statusCode, w)
}
