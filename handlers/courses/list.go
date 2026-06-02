package courses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mobile-server/database"
	"mobile-server/handlers"
)

// ListHandler retorna todos os cursos com informacao de inscricao do usuario.
func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	user, ok := handlers.RequireAuthenticatedUser(w, username, password, writeListError)
	if !ok {
		return
	}

	courses, err := database.GetAllCoursesWithEnrollment(user.ID)
	if err != nil {
		fmt.Printf("Erro ao buscar cursos: %v\n", err)
		http.Error(w, "Erro ao buscar cursos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func writeListError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
