package courses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mobile-server/database"
	"mobile-server/handlers"
	"mobile-server/models"
)

// CreateHandler cria um novo curso somente para usuarios admins.
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeCreateCourseError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	user, ok := handlers.RequireAuthenticatedUser(w, req.Username, req.Password, writeCreateCourseError)
	if !ok {
		return
	}

	if !user.IsAdmin {
		writeCreateCourseError(w, http.StatusForbidden, "Apenas administradores podem criar cursos")
		return
	}

	if req.Name == "" || req.DayOfWeek == "" || req.Time == "" || req.Location == "" || req.Professor == "" {
		writeCreateCourseError(w, http.StatusBadRequest, "Campos obrigatorios: username, password, name, day_of_week, time, location, professor")
		return
	}

	courseID, err := database.CreateCourse(req)
	if err != nil {
		fmt.Printf("Erro ao criar curso: %v\n", err)
		writeCreateCourseError(w, http.StatusInternalServerError, "Falha ao criar curso")
		return
	}

	resp := models.CreateCourseResponse{
		Success:  true,
		Message:  "Curso criado com sucesso",
		CourseID: &courseID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func writeCreateCourseError(w http.ResponseWriter, status int, message string) {
	resp := models.CreateCourseResponse{
		Success: false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
