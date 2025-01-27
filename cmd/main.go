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

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning : .env file not found %v", err)
	}

	// initiate database
	database.ConnectDatabase()

	// create route api
	router := mux.NewRouter()

	// authentication
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/refresh", handlers.GetRefreshToken).Methods("POST")
	router.HandleFunc("/api/register/seeker", handlers.SeekerRegister).Methods("POST")
	router.HandleFunc("/api/register/employer", handlers.EmployerRegister).Methods("POST")

	// public access
	router.HandleFunc("/api/jobs", handlers.GetAllJobs).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", handlers.GetJobByID).Methods("GET")

	// Protected endpoint with middleware
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/api/me", handlers.AuthMe).Methods("POST")

	// seeker
	protected.HandleFunc("/api/jobs/{id}/apply", handlers.ApplyToJob).Methods("POST")
	protected.HandleFunc("/api/seeker/profile", handlers.GetUserSeekerProfile).Methods("GET")
	protected.HandleFunc("/api/seeker/profile", handlers.UpdateUserSeekerProfile).Methods("PUT")
	protected.HandleFunc("/api/seeker/applications", handlers.GetSeekerJobApplication).Methods("GET")
	protected.HandleFunc("/api/seeker/profile/experience", handlers.AddUserSeekerExperience).Methods("POST")
	protected.HandleFunc("/api/seeker/profile/experience/{id}", handlers.UpdateUserSeekerExperience).Methods("PUT")

	// next feature
	// protected.HandleFunc("/api/seeker/applications/{id}", handlers.CancelSeekerJobApplication).Methods("DELETE")

	// employer
	protected.HandleFunc("/api/employer/jobs", handlers.CreateJob).Methods("POST")
	protected.HandleFunc("/api/employer/jobs/{id}", handlers.UpdateJob).Methods("PUT")
	protected.HandleFunc("/api/employer/jobs/{id}", handlers.DeleteJob).Methods("DELETE")
	protected.HandleFunc("/api/employer/jobs", handlers.GetAllEmployerPostedJobs).Methods("GET")
	protected.HandleFunc("/api/employer/jobs/{id}/applications", handlers.UpdateApplicationStatus).Methods("PUT")
	protected.HandleFunc("/api/employer/jobs/{id}/applications", handlers.GetEmployerJobApplications).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
