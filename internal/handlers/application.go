package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// ApplyToJob handles job application by job seekers.
func ApplyToJob(w http.ResponseWriter, r *http.Request) {
	// Check if the user is a job seeker
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "job_seeker" {
		http.Error(w, "Unauthorized: Job seeker only", http.StatusUnauthorized)
		return
	}

	jobID := mux.Vars(r)["id"]
	var job models.Job
	if err := database.DB.First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Create application
	application := models.Application{
		JobID:  job.ID,
		UserID: uint(claims["user_id"].(float64)),
		Status: "Pending",
	}

	if err := database.DB.Create(&application).Error; err != nil {
		http.Error(w, "Failed to apply to job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Application submitted successfully"})
}


// GetApplicationsByJobID handles getting all applications for a specific job (employer only).
func GetApplicationsByJobID(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	jobID := mux.Vars(r)["id"]
	var job models.Job
	if err := database.DB.First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Fetch applications for the job
	var applications []models.Application
	if err := database.DB.Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}
// GetApplicationsByUserID handles getting all applications for a specific user (job seeker only).
func GetApplicationsByUserID(w http.ResponseWriter, r *http.Request) {
	// Check if the user is a job seeker
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "job_seeker" {
		http.Error(w, "Unauthorized: Job seeker only", http.StatusUnauthorized)
		return
	}

	// Extract user ID from JWT claims
	userID := uint(claims["user_id"].(float64))

	// Fetch applications for the user
	var applications []models.Application
	if err := database.DB.Where("user_id = ?", userID).Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}


// UpdateApplicationStatus handles updating the status of a job application (employer only).
func UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	applicationID := mux.Vars(r)["id"]
	var application models.Application
	if err := database.DB.First(&application, applicationID).Error; err != nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	// Ensure the application belongs to one of the employer's jobs
	userID := claims["user_id"].(float64)
	var company models.Company
	if err := database.DB.Where("user_id = ?", uint(userID)).First(&company).Error; err != nil {
		http.Error(w, "Failed to find employer's company", http.StatusInternalServerError)
		return
	}

	var job models.Job
	if err := database.DB.First(&job, application.JobID).Error; err != nil || job.CompanyID != company.ID {
		http.Error(w, "Unauthorized: Cannot update applications for jobs you do not own", http.StatusForbidden)
		return
	}

	// Update the status of the application
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	application.Status = payload.Status
	if err := database.DB.Save(&application).Error; err != nil {
		http.Error(w, "Failed to update application status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Application status updated successfully"})
}
