package config

import (
	"net/http"

	"github.com/ddomd/feeder/utils/http_helpers"
)

func HandleServerReady (write http.ResponseWriter, req *http.Request) {
	http_helpers.RespondWithJson(write, http.StatusOK, map[string]string{"status": "ok"})
}

func HandleServerError (write http.ResponseWriter, req *http.Request) {
	http_helpers.RespondWithError(write, http.StatusInternalServerError, "Internal Server Error")
}
