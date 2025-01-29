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
	"github.com/rs/cors"
)

func main() {
	// Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning : .env file not found %v", err)
	}

	// Initiate database
	database.ConnectDatabase()

	// Create router
	router := mux.NewRouter()

	// Middleware API Key
	// router.Use(middleware.APIKeyMiddleware)

	// Authentication Routes
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/logout", handlers.Logout).Methods("POST")
	router.HandleFunc("/api/refresh", handlers.GetRefreshToken).Methods("POST")
	router.HandleFunc("/api/register/seeker", handlers.SeekerRegister).Methods("POST")
	router.HandleFunc("/api/register/employer", handlers.EmployerRegister).Methods("POST")

	// Public Routes
	router.HandleFunc("/api/jobs", handlers.GetAllJobs).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", handlers.GetJobByID).Methods("GET")
	router.HandleFunc("/api/employer/{id}", handlers.GetEmployerCompanyProfile).Methods("GET")

	// Protected Routes with JWT Middleware
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTMiddleware)
	protected.HandleFunc("/api/me", handlers.AuthMe).Methods("POST")

	// Seeker Routes
	protected.HandleFunc("/api/jobs/{id}/apply", handlers.ApplyToJob).Methods("POST")
	protected.HandleFunc("/api/seeker/profile", handlers.GetSeekerProfile).Methods("GET")
	protected.HandleFunc("/api/seeker/profile", handlers.UpdateSeekerProfile).Methods("PUT")
	protected.HandleFunc("/api/seeker/applications", handlers.GetSeekerJobApplications).Methods("GET")
	protected.HandleFunc("/api/seeker/profile/experience", handlers.AddUserSeekerExperience).Methods("POST")
	protected.HandleFunc("/api/seeker/profile/experience/{id}", handlers.UpdateUserSeekerExperience).Methods("PUT")

	// Employer Routes
	protected.HandleFunc("/api/employer/jobs", handlers.CreateJob).Methods("POST")
	protected.HandleFunc("/api/employer/jobs/{id}", handlers.UpdateJob).Methods("PUT")
	protected.HandleFunc("/api/employer/jobs/{id}", handlers.DeleteJob).Methods("DELETE")
	protected.HandleFunc("/api/employer/jobs", handlers.GetAllEmployerPostedJobs).Methods("GET")
	protected.HandleFunc("/api/employer/profile", handlers.GetEmployerProfile).Methods("GET")
	protected.HandleFunc("/api/employer/profile", handlers.UpdateEmployerProfile).Methods("PUT")
	protected.HandleFunc("/api/employer/jobs/{id}/applications", handlers.UpdateApplicationStatus).Methods("PUT")
	protected.HandleFunc("/api/employer/jobs/{id}/applications", handlers.GetEmployerJobApplications).Methods("GET")

	client_host := os.Getenv("CLIENT_HOST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{client_host},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-API-key"},
		AllowCredentials: true,
	})

	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := corsHandler.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
