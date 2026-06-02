package models

// CourseWithEnrollment representa um curso com informação de inscrição do usuário.
type CourseWithEnrollment struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	DayOfWeek    string `json:"day_of_week"`
	Time         string `json:"time"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	Professor    string `json:"professor"`
	IsEnrolled   bool   `json:"is_enrolled"`
	EnrollmentID *int64 `json:"enrollment_id,omitempty"`
}
