package config

import (
	"net/http"

	"github.com/ddomd/feeder/internal/database"
)

type ApiConfig struct {
	Db *database.Queries
	Port string
}

type AuthenticatedHandler func(http.ResponseWriter, *http.Request, database.User)
