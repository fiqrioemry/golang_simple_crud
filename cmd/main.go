package main

import (
	"golang_project/internal/database"
	"golang_project/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load environtment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning : .env file not found %v", err)
	}

	// Connect to the database
	database.ConnectDatabase()

	// create new router
	router := mux.NewRouter()


	// define API endpoint
	router.HandleFunc("/auth/register/seeker", handlers.SeekerRegister).Methods("POST")
	router.HandleFunc("/auth/register/employer", handlers.EmployerRegister).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	
	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
