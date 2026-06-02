package models

// Course representa um curso com informações básicas.
type Course struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DayOfWeek   string `json:"day_of_week"` // e.g., "Monday" or number
	Time        string `json:"time"`        // e.g., "18:30"
	Description string `json:"description"`
	Location    string `json:"location"`
	Professor   string `json:"professor"`
}

// CourseDate representa uma data específica associada a um curso.
type CourseDate struct {
	ID       int64  `json:"id"`
	CourseID int64  `json:"course_id"`
	Date     string `json:"date"` // ISO date: YYYY-MM-DD
}

// CreateCourseRequest representa uma requisição de criação de curso.
type CreateCourseRequest struct {
	Name        string `json:"name"`
	DayOfWeek   string `json:"day_of_week"`
	Time        string `json:"time"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Professor   string `json:"professor"`
}

// CreateCourseResponse representa a resposta de criação de curso.
type CreateCourseResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	CourseID *int64 `json:"course_id,omitempty"`
}
