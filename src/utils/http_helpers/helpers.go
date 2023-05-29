package http_helpers

import (
	"encoding/json"
	"net/http"
)

func RespondWithJson(write http.ResponseWriter, statusCode int, payload interface{}) {
	write.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload); if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
		return
	}

	write.WriteHeader(statusCode)
	write.Write(response)
}

func RespondWithError(write http.ResponseWriter, statusCode int, errorMsg string) {
	err := struct {
		Error string `json:"error"`
	}{
		errorMsg,
	}

	RespondWithJson(write, statusCode, err)
}