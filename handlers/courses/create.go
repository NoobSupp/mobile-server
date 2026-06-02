package courses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mobile-server/database"
	"mobile-server/models"
)

// CreateHandler cria um novo curso (somente admins).
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extract user_id from query string to verify admin status
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "Parâmetro user_id é obrigatório",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "user_id inválido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Verify user is admin
	user, err := database.GetUserByID(userID)
	if err != nil {
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "Usuário não encontrado",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !user.IsAdmin {
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "Apenas administradores podem criar cursos",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req models.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "JSON inválido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if req.Name == "" || req.DayOfWeek == "" || req.Time == "" || req.Location == "" || req.Professor == "" {
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "Campos obrigatórios: name, day_of_week, time, location, professor",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	courseID, err := database.CreateCourse(req)
	if err != nil {
		fmt.Printf("Erro ao criar curso: %v\n", err)
		resp := models.CreateCourseResponse{
			Success: false,
			Message: "Falha ao criar curso",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
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
