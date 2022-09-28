package util

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(payload any, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	encoded, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"Msg": "Failed to marshal response", "Error": "See server log for more information"}`))
		return
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(encoded)
}
