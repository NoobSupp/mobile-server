package courses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mobile-server/database"
	"mobile-server/handlers"
	"mobile-server/models"
)

// EnrollHandler inscreve um usuario em um curso.
func EnrollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.EnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	user, ok := handlers.RequireAuthenticatedUser(w, req.Username, req.Password, writeEnrollmentError)
	if !ok {
		return
	}

	if req.CourseID == 0 {
		writeEnrollmentError(w, http.StatusBadRequest, "course_id e obrigatorio")
		return
	}

	enrollmentID, err := database.CreateEnrollment(user.ID, req.CourseID)
	if err != nil {
		fmt.Printf("Erro ao inscrever usuario: %v\n", err)
		writeEnrollmentError(w, http.StatusConflict, "Falha ao inscrever no curso (usuario ja pode estar inscrito)")
		return
	}

	resp := models.EnrollmentResponse{
		Success:      true,
		Message:      "Inscricao realizada com sucesso",
		EnrollmentID: &enrollmentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func writeEnrollmentError(w http.ResponseWriter, status int, message string) {
	resp := models.EnrollmentResponse{
		Success: false,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
