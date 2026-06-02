package models

// Enrollment representa uma inscrição de um usuário em um curso.
type Enrollment struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"user_id"`
	CourseID int64 `json:"course_id"`
}

// EnrollmentRequest representa uma requisição de inscrição em um curso.
type EnrollmentRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	CourseID int64  `json:"course_id"`
}

// EnrollmentResponse representa a resposta de uma inscrição.
type EnrollmentResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	EnrollmentID *int64 `json:"enrollment_id,omitempty"`
}
