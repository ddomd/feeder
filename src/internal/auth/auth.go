package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAuthHeader(header http.Header) (string, error){
	authHeader := header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Empty authentication header")
	}

	authFields := strings.Split(authHeader, " ")

	if len(authFields) < 2 || authFields[0] != "ApiKey" {
		return "", errors.New("Malformed authorization header")
	}

	return authFields[1], nil
}