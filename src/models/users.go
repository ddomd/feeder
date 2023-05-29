package models

import (
	"time"

	"github.com/ddomd/feeder/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username string `json:"username"`
	APIKey string `json:"api_key"`
}

func ConvertToUserModel(user database.User) User{
	return User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username: user.Username,
		APIKey: user.ApiKey,
	}
}
