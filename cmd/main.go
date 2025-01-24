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


func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning : .env file not found %v", err)
	}

	database.ConnectDatabase()

	router 		:= mux.NewRouter().PathPrefix("/api")
	auth 		:= router.PathPrefix("/auth").Subrouter()
	protected 	:= router.PathPrefix("").Subrouter()

	auth.HandleFunc("/register/seeker", handlers.SeekerRegister).Methods("POST")
	auth.HandleFunc("/register/employer", handlers.EmployerRegister).Methods("POST")
	auth.HandleFunc("/login", handlers.Login).Methods("POST")
	auth.protected.HandleFunc("/me", handlers.AuthMe).Methods("POST")

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	 	if err := http.ListenAndServe(":"+port, router); err != nil {
	 		log.Fatalf("Failed to start server: %v", err)
	}

}