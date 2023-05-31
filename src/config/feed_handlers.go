package config

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/ddomd/feeder/models"
	"github.com/ddomd/feeder/utils/http_helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleCreateFeed(write http.ResponseWriter, req *http.Request, user database.User) {
	type Parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	params := Parameters{}

	err := json.NewDecoder(req.Body).Decode(&params); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	feed, err := cfg.Db.CreateFeed(req.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	feed_follow, err := cfg.Db.AddFollow(req.Context(), database.AddFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	res := struct{
		Feed       models.Feed `json:"feed"`
		FeedFollow models.FeedFollow `json:"feed_follow"`
	}{
		models.ConvertToFeedModel(feed),
		models.ConvertToFeedFollowModel(feed_follow),
	}

	http_helpers.RespondWithJson(write, http.StatusAccepted, res)
}

func(cfg *ApiConfig) HandleDeleteFeed(write http.ResponseWriter, req *http.Request, user database.User) {
	feedId, err := uuid.Parse(chi.URLParam(req, "feedId")); if err != nil {
		http_helpers.RespondWithError(write, http.StatusBadRequest, err.Error())
		return
	}

	err = cfg.Db.DeleteFeed(req.Context(), database.DeleteFeedParams{
		ID:     feedId,
		UserID: user.ID,
		}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusOK, struct{}{})
}

func (cfg *ApiConfig) HandleGetAllFeeds(write http.ResponseWriter, req *http.Request) {
	feeds, err := cfg.Db.GetAllFeeds(req.Context()); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusOK, models.ConvertToFeedsModel(feeds))
}
