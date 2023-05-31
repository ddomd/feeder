package config

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/ddomd/feeder/utils/http_helpers"
	"github.com/ddomd/feeder/models"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleCreateUser(write http.ResponseWriter, req *http.Request) {
	type Parameters struct {
		Username string `json:"username"`
	}

	params := Parameters{}

	err := json.NewDecoder(req.Body).Decode(&params); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.Db.CreateUser(req.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Username:  params.Username,
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusOK, models.ConvertToUserModel(user))
}

func (cfg *ApiConfig) HandleGetUserByAPIKey(write http.ResponseWriter, req *http.Request, user database.User) {
	http_helpers.RespondWithJson(write, http.StatusAccepted, models.ConvertToUserModel(user))
}
