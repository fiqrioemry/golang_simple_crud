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

	router := mux.NewRouter()

	router.HandleFunc("/register/seeker", handlers.SeekerRegister).Methods("POST")
	router.HandleFunc("/register/employer", handlers.EmployerRegister).Methods("POST")
	router.HandleFunc("/login", handlers.EmployerRegister).Methods("POST")

	port := os.Getenv()
	if port == ""{
		port = "5000"
	}

	log.Printf("Starting server on :%s...", port)
	 	if err := http.ListenAndServe(":"+port, router); err != nil {
	 		log.Fatalf("Failed to start server: %v", err)
	}
	

}

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning : .env file not found %v", err)
	}

	database.ConnectDatabase()

	router := mux.NewRouter()

	router.HandleFunc("/register/seeker", handlers.SeekerRegister).Methods("POST")
	router.HandleFunc("/register/employer", handlers.EmployerRegister).Methods("POST")
	router.HandleFunc("/login", handlers.EmployerRegister).Methods("POST")

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	 	if err := http.ListenAndServe(":"+port, router); err != nil {
	 		log.Fatalf("Failed to start server: %v", err)
	}
	

}