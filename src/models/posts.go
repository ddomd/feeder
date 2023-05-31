package models

import (
	"database/sql"
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description string     `json:"description"`
	PublishedAt time.Time  `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

func convertFromNullTime(nulltime sql.NullTime) time.Time{
	if nulltime.Valid {
		return nulltime.Time
	}
	return time.Time{}
}

func convertFromNullString(nullstr sql.NullString) string{
	if nullstr.Valid {
		return nullstr.String
	}
	return ""
}

func ConvertToPostModel(post database.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: convertFromNullString(post.Description),
		PublishedAt: convertFromNullTime(post.PublishedAt),
		FeedID:      post.FeedID,
	}
}

func ConvertToPostsModel(posts []database.Post) Posts {
	postSlice := make([]Post, len(posts))

	for i, _ := range posts {
		postSlice[i] = ConvertToPostModel(posts[i])
	}

	return Posts{Posts: postSlice}
}