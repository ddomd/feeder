package config

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/ddomd/feeder/models"
	"github.com/ddomd/feeder/utils/http_helpers"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleCreateFeed(write http.ResponseWriter, req *http.Request, user database.User) {
	type Parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}

	params := Parameters{}

	err := json.NewDecoder(req.Body).Decode(&params); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	feed, err := cfg.Db.CreateFeed(req.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusAccepted, models.ConvertToFeedModel(feed))
}

func (cfg *ApiConfig) HandleGetAllFeeds(write http.ResponseWriter, req *http.Request) {
	feeds, err := cfg.Db.GetAllFeeds(req.Context()); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusOK, models.ConvertToFeedsModel(feeds))
}