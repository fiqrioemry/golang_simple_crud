package main

import (
	"golang_project/internal/database"
	"golang_project/internal/handlers"
	"golang_project/internal/middleware"
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

	router	:= mux.NewRouter()

	router.HandleFunc("/api/register/seeker", handlers.SeekerRegister).Methods("POST")
	router.HandleFunc("/api/register/employer", handlers.EmployerRegister).Methods("POST")
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/jobs", handlers.GetAllJobs).Methods("GET") 
	router.HandleFunc("/api/jobs/{id}", handlers.GetJobByID).Methods("GET")  
	router.HandleFunc("/refresh", handlers.GetRefreshToken).Methods("POST")
	// router.HandleFunc("/api/jobs/employer/{id}", handlers.GetAllEmployerJobs).Methods("GET") 


	// Protected endpoint with middleware
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/me", handlers.AuthMe).Methods("POST")



	// seeker
	protected.HandleFunc("/api/jobs/{id}/apply", handlers.ApplyToJob).Methods("POST") 
	protected.HandleFunc("/api/applications/user/{id}", handlers.GetApplicationsByUserID).Methods("GET") 
	protected.HandleFunc("/api/seeker/profile", handlers.GetUserSeekerProfile).Methods("GET")   // Get profile
	protected.HandleFunc("/api/seeker/profile", handlers.UpdateUserSeekerProfile).Methods("PUT") // Update profile
	


	// employer
	protected.HandleFunc("/api/jobs/{id}", handlers.UpdateJob).Methods("PUT") 
	protected.HandleFunc("/api/jobs", handlers.CreateJob).Methods("POST")  
	protected.HandleFunc("/api/jobs/{id}", handlers.DeleteJob).Methods("DELETE") 
	protected.HandleFunc("/api/applications/job/{id}", handlers.GetApplicationsByJobID).Methods("GET") 
	protected.HandleFunc("/api/applications/status", handlers.UpdateApplicationStatus).Methods("PUT")


	// protected.HandleFunc("/api/user/employer", handlers.GetUserEmployerProfile).Methods("GET")  
	// protected.HandleFunc("/api/user/employer", handlers.EditUserEmployerProfile).Methods("PUT")  
	// router.HandleFunc("/api/jobs/employer", handlers.GetAllEmployerPostedJobs).Methods("GET") 



	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	 	if err := http.ListenAndServe(":"+port, router); err != nil {
	 		log.Fatalf("Failed to start server: %v", err)
	}

}