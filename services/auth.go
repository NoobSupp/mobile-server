package services

import (
	"database/sql"
	"fmt"

	"mobile-server/database"
	"mobile-server/models"

	"golang.org/x/crypto/bcrypt"
)

func Authenticate(req models.LoginRequest) models.LoginResponse {
	if _, err := AuthenticateUser(req.Username, req.Password); err != nil {
		return models.LoginResponse{Success: false, Message: err.Error()}
	}

	return models.LoginResponse{
		Success: true,
		Message: "Login realizado com sucesso",
	}
}

func AuthenticateUser(username, password string) (models.User, error) {
	user, err := database.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("Usuario ou senha invalidos")
		}
		return models.User{}, fmt.Errorf("Erro interno ao consultar usuario")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, fmt.Errorf("Usuario ou senha invalidos")
	}

	return user, nil
}
