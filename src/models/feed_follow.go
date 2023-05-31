package models

import (
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/google/uuid"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type FeedFollows struct {
	Follows []FeedFollow `json:"feed_follows"`
}

func ConvertToFeedFollowModel(feed database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		UserID:    feed.UserID,
		FeedID:    feed.FeedID,
	}
}

func ConvertToFeedFollowsModel(feeds []database.FeedFollow) FeedFollows {
	feedFollows := make([]FeedFollow, len(feeds))

	for i, feedFollow := range feeds {
		feedFollows[i] = ConvertToFeedFollowModel(feedFollow)
	}

	return FeedFollows{feedFollows}
}
