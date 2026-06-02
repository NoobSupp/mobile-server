package services

import (
	"database/sql"

	"mobile-server/database"
	"mobile-server/models"

	"golang.org/x/crypto/bcrypt"
)

func Authenticate(req models.LoginRequest) models.LoginResponse {
	user, err := database.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.LoginResponse{Success: false, Message: "Usuário ou senha inválidos"}
		}
		return models.LoginResponse{Success: false, Message: "Erro interno ao consultar usuário"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return models.LoginResponse{Success: false, Message: "Usuário ou senha inválidos"}
	}

	return models.LoginResponse{
		Success: true,
		Message: "Login realizado com sucesso",
	}
}
