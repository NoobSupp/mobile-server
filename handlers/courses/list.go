package courses

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"mobile-server/database"
)

// ListHandler retorna todos os cursos com informação de inscrição do usuário.
func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair username da query string
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Parâmetro username é obrigatório", http.StatusBadRequest)
		return
	}

	// Buscar usuário pelo username
	user, err := database.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}
		fmt.Printf("Erro ao buscar usuário: %v\n", err)
		http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
		return
	}

	// Buscar cursos com informação de inscrição do usuário
	courses, err := database.GetAllCoursesWithEnrollment(user.ID)
	if err != nil {
		fmt.Printf("Erro ao buscar cursos: %v\n", err)
		http.Error(w, "Erro ao buscar cursos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
