package handlers

import (
	"net/http"

	"mobile-server/models"
	"mobile-server/services"
)

type ErrorWriter func(w http.ResponseWriter, status int, message string)

func RequireAuthenticatedUser(w http.ResponseWriter, username, password string, writeError ErrorWriter) (models.User, bool) {
	if username == "" || password == "" {
		writeError(w, http.StatusBadRequest, "Campos username e password sao obrigatorios")
		return models.User{}, false
	}

	user, err := services.AuthenticateUser(username, password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return models.User{}, false
	}

	return user, true
}
