package config

import (
	"net/http"

	"github.com/ddomd/feeder/internal/auth"
	"github.com/ddomd/feeder/utils/http_helpers"
)

func (cfg *ApiConfig) AuthMiddleware(handler AuthenticatedHandler) http.HandlerFunc {
	return func(write http.ResponseWriter, req *http.Request) {
		key, err := auth.GetAuthHeader(req.Header); if err != nil {
			http_helpers.RespondWithError(write, http.StatusUnauthorized, err.Error())
			return
		}

		authedUser, err := cfg.Db.GetUserByAPIKey(req.Context(), key); if err != nil {
			http_helpers.RespondWithError(write, http.StatusUnauthorized, err.Error())
			return
		}

		handler(write, req, authedUser)
	}
}