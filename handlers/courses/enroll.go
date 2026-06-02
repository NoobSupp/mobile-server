package courses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mobile-server/database"
	"mobile-server/models"
)

// EnrollHandler inscreve um usuário em um curso.
func EnrollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.EnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.CourseID == 0 {
		resp := models.EnrollmentResponse{
			Success: false,
			Message: "user_id e course_id são obrigatórios",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	enrollmentID, err := database.CreateEnrollment(req.UserID, req.CourseID)
	if err != nil {
		fmt.Printf("Erro ao inscrever usuário: %v\n", err)
		resp := models.EnrollmentResponse{
			Success: false,
			Message: "Falha ao inscrever no curso (usuário já pode estar inscrito)",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.EnrollmentResponse{
		Success:      true,
		Message:      "Inscrição realizada com sucesso",
		EnrollmentID: &enrollmentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
