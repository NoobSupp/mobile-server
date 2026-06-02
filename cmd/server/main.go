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

	http.HandleFunc("/", handlers.LocalNetworkOnly(healthHandler))
	http.HandleFunc("/login", handlers.LocalNetworkOnly(handlers.LoginHandler))
	http.HandleFunc("/cursos", handlers.LocalNetworkOnly(courses.Router))
	http.HandleFunc("/cursos/inscrever", handlers.LocalNetworkOnly(courses.EnrollHandler))

	fmt.Println("Servidor iniciado na porta 4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
