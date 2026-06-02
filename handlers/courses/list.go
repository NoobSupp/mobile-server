package courses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mobile-server/database"
)

// ListHandler retorna todos os cursos com informação de inscrição do usuário.
func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair user_id da query string
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Parâmetro user_id é obrigatório", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "user_id inválido", http.StatusBadRequest)
		return
	}

	courses, err := database.GetAllCoursesWithEnrollment(userID)
	if err != nil {
		fmt.Printf("Erro ao buscar cursos: %v\n", err)
		http.Error(w, "Erro ao buscar cursos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
