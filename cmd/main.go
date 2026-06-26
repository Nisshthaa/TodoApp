package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/handlers"
)

func main() {
	if err := database.OpenConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	); err != nil {
		log.Panicf("Failed to initialize and migrate database with error: %v", err)
	}

	defer database.Todo.Close()

	router := chi.NewRouter()

	router.Route("/v1", func(v1 chi.Router) {
		v1.Post("/register", handlers.RegisterUser)

	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Panicf("Failed to start server with error: %v", err)

	}
	log.Println("Server started at :8080")

}
