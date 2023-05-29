package config

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/ddomd/feeder/models"
	"github.com/ddomd/feeder/utils/http_helpers"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
)


func (cfg *ApiConfig) HandleAddFollow(write http.ResponseWriter, req *http.Request, user database.User) {
	type Parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := Parameters{}
	err := json.NewDecoder(req.Body).Decode(&params); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	follow, err := cfg.Db.AddFollow(req.Context(), database.AddFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusAccepted, models.ConvertToFeedFollowModel(follow))
}

func (cfg *ApiConfig) HandleGetUserFollows(write http.ResponseWriter, req *http.Request, user database.User) {
	follows, err := cfg.Db.GetAllUserFollows(req.Context(), user.ID); if err != nil {
		http_helpers.RespondWithError(write, http.StatusUnauthorized, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusAccepted, models.ConvertToFeedFollowsModel(follows))
}

func (cfg *ApiConfig) HandleRemoveFollow(write http.ResponseWriter, req *http.Request, user database.User) {
	idParam, err := uuid.Parse(chi.URLParam(req, "followId")); if err != nil {
		http_helpers.RespondWithError(write, http.StatusBadRequest, err.Error())
		return
	}

	err = cfg.Db.RemoveFollow(req.Context(), database.RemoveFollowParams{
		ID: idParam,
		UserID: user.ID,
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, "Couldn't delete feed follow")
		return
	}

	http_helpers.RespondWithJson(write, http.StatusOK, struct{}{})
}
