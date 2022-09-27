package util

import (
	"encoding/json"
	"net/http"
)

func WriteErrorResponse(msg string, err error, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	encoded, marshalErr := json.Marshal(
		struct {
			Msg   string
			Error error `json:",omitempty"`
		}{
			Msg:   msg,
			Error: err,
		},
	)
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"Msg": "Failed to marshal error message", "Error": "See server log for more information"}`))
		return
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(encoded)
}
