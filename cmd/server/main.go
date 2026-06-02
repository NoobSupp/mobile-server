package main

import (
	"fmt"
	"log"
	"net/http"

	"mobile-server/database"
	"mobile-server/handlers"
	"mobile-server/handlers/courses"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Servidor funcionando"))
}

func main() {
	if err := database.Init("data/server.db"); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/cursos", courses.Router)
	http.HandleFunc("/cursos/inscrever", courses.EnrollHandler)

	fmt.Println("Servidor iniciado na porta 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
