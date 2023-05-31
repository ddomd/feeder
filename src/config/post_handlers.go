package config

import (
	"net/http"

	"github.com/ddomd/feeder/internal/database"
	"github.com/ddomd/feeder/models"
	"github.com/ddomd/feeder/utils/http_helpers"
)

func (cfg *ApiConfig) HandleGetUserPosts(write http.ResponseWriter, req *http.Request, user database.User) {
	posts, err := cfg.Db.GetUserPosts(req.Context(), database.GetUserPostsParams{
		UserID: user.ID,
		Limit:  int32(10),
	}); if err != nil {
		http_helpers.RespondWithError(write, http.StatusInternalServerError, err.Error())
		return
	}

	http_helpers.RespondWithJson(write, http.StatusOK, models.ConvertToPostsModel(posts))
}