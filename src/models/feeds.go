package models

import (
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string `json:"name"`
	URL string `json:"url"`
	UserID uuid.UUID `json:"user_id"`
}

type Feeds struct {
	Feeds []Feed `json:"feeds"`
}

func ConvertToFeedModel(feed database.Feed) Feed{
	return Feed{
		ID: feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name: feed.Name,
		URL: feed.Url,
		UserID: feed.UserID,
	}
}

func ConvertToFeedsModel(feeds []database.Feed) Feeds{
	convertedFeeds := Feeds{}
	feedSlice := make([]Feed, len(feeds))

	for i, feed := range feeds {
		feedSlice[i] = ConvertToFeedModel(feed)
	}

	convertedFeeds.Feeds = feedSlice

	return convertedFeeds
}
