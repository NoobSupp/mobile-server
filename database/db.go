package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"mobile-server/models"

	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Init(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("falha ao criar pasta do banco: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("falha ao abrir banco de dados: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("falha ao conectar no banco de dados: %w", err)
	}

	DB = db
	return migrate()
}

func Close() error {
	if DB == nil {
		return nil
	}
	return DB.Close()
}

func migrate() error {
	// Enable foreign key constraints
	if _, err := DB.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("falha ao habilitar foreign_keys: %w", err)
	}

	const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_admin INTEGER DEFAULT 0
);
`

	if _, err := DB.Exec(createUsersTable); err != nil {
		return fmt.Errorf("falha ao criar tabela users: %w", err)
	}

	// Add is_admin column if it doesn't exist (for existing databases)
	if _, err := DB.Exec(`ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT 0;`); err != nil {
		// Ignore error if column already exists
	}

	const createCoursesTable = `
CREATE TABLE IF NOT EXISTS courses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    day_of_week TEXT NOT NULL,
    time TEXT NOT NULL,
    description TEXT,
    location TEXT,
    professor TEXT
);
`
	if _, err := DB.Exec(createCoursesTable); err != nil {
		return fmt.Errorf("falha ao criar tabela courses: %w", err)
	}

	const createCourseDatesTable = `
CREATE TABLE IF NOT EXISTS course_dates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    course_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    FOREIGN KEY(course_id) REFERENCES courses(id) ON DELETE CASCADE
);
`
	if _, err := DB.Exec(createCourseDatesTable); err != nil {
		return fmt.Errorf("falha ao criar tabela course_dates: %w", err)
	}

	const createEnrollmentsTable = `
CREATE TABLE IF NOT EXISTS enrollments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    UNIQUE(user_id, course_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(course_id) REFERENCES courses(id) ON DELETE CASCADE
);
`
	if _, err := DB.Exec(createEnrollmentsTable); err != nil {
		return fmt.Errorf("falha ao criar tabela enrollments: %w", err)
	}

	return ensureAdminAccount()
}

func ensureAdminAccount() error {
	const adminUsername = "admin"
	if _, err := GetUserByUsername(adminUsername); err == nil {
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("falha ao gerar hash da senha padrão: %w", err)
	}

	_, err = DB.Exec(`INSERT OR IGNORE INTO users (username, password) VALUES (?, ?)`, adminUsername, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("falha ao inserir usuário admin: %w", err)
	}

	return nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	var isAdmin int
	row := DB.QueryRow(`SELECT id, username, password, is_admin FROM users WHERE username = ?`, username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &isAdmin); err != nil {
		return models.User{}, err
	}
	user.IsAdmin = isAdmin != 0
	return user, nil
}

// GetUserByID retorna um usuário pelo ID.
func GetUserByID(userID int64) (models.User, error) {
	var user models.User
	var isAdmin int
	row := DB.QueryRow(`SELECT id, username, password, is_admin FROM users WHERE id = ?`, userID)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &isAdmin); err != nil {
		return models.User{}, err
	}
	user.IsAdmin = isAdmin != 0
	return user, nil
}

// GetAllCoursesWithEnrollment retorna todos os cursos com informação de inscrição do usuário.
func GetAllCoursesWithEnrollment(userID int64) ([]models.CourseWithEnrollment, error) {
	query := `
SELECT c.id, c.name, c.day_of_week, c.time, c.description, c.location, c.professor,
       COALESCE(e.id, 0) AS enrollment_id
FROM courses c
LEFT JOIN enrollments e ON c.id = e.course_id AND e.user_id = ?
ORDER BY c.id
`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar cursos: %w", err)
	}
	defer rows.Close()

	var courses []models.CourseWithEnrollment
	for rows.Next() {
		var course models.CourseWithEnrollment
		var enrollmentID int64
		if err := rows.Scan(&course.ID, &course.Name, &course.DayOfWeek, &course.Time,
			&course.Description, &course.Location, &course.Professor, &enrollmentID); err != nil {
			return nil, fmt.Errorf("falha ao scanear curso: %w", err)
		}
		course.IsEnrolled = enrollmentID > 0
		if enrollmentID > 0 {
			course.EnrollmentID = &enrollmentID
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar cursos: %w", err)
	}

	return courses, nil
}

// CreateEnrollment cria uma nova inscrição de usuário em um curso.
func CreateEnrollment(userID, courseID int64) (int64, error) {
	result, err := DB.Exec(
		`INSERT INTO enrollments (user_id, course_id) VALUES (?, ?)`,
		userID, courseID,
	)
	if err != nil {
		return 0, fmt.Errorf("falha ao criar inscrição: %w", err)
	}

	enrollmentID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("falha ao obter ID da inscrição: %w", err)
	}

	return enrollmentID, nil
}

// CreateCourse cria um novo curso no banco de dados.
func CreateCourse(req models.CreateCourseRequest) (int64, error) {
	result, err := DB.Exec(
		`INSERT INTO courses (name, day_of_week, time, description, location, professor) VALUES (?, ?, ?, ?, ?, ?)`,
		req.Name, req.DayOfWeek, req.Time, req.Description, req.Location, req.Professor,
	)
	if err != nil {
		return 0, fmt.Errorf("falha ao criar curso: %w", err)
	}

	courseID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("falha ao obter ID do curso: %w", err)
	}

	return courseID, nil
}
