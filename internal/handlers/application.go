package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

func ApplyToJob(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
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


func GetApplicationsByJobID(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Extract job ID from the request URL
	jobID := mux.Vars(r)["id"]

	// Verify the job exists
	var job models.Job
	if err := database.DB.First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Fetch applications for the job and preload related user and profile data
	var applications []models.Application
	if err := database.DB.Preload("User.Profile").Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	// If no applications are found, return a 404
	if len(applications) == 0 {
		http.Error(w, "No applications found for this job", http.StatusNotFound)
		return
	}

	// Transform the response to include user and profile details
	var responses []map[string]interface{}
	for _, application := range applications {
		response := map[string]interface{}{
			"id"			:application.ID,
			"status"		:application.Status,
			"created_at"	:application.CreatedAt,
			"updated_at"	:application.UpdatedAt,
			"user"			:map[string]interface{}{
				"id"		:application.UserID,
				"name"		:application.User.Name, 
				"email"		:application.User.Email, 
				"profile"	:map[string]interface{}{
					"bio"	:application.User.Profile.Bio,    
					"resume":application.User.Profile.Resume, 
					"skills":application.User.Profile.Skills, 
				},
			},
		}
		responses = append(responses, response)
	}

	// Return the transformed response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// GetApplicationsByUserID handles getting all applications for a specific user (job seeker only).
func GetApplicationsByUserID(w http.ResponseWriter, r *http.Request) {
	// Check if the user is a job seeker
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Job seeker only", http.StatusUnauthorized)
		return
	}

	// Extract user ID from JWT claims
	userID := uint(claims["user_id"].(float64))

	// Fetch applications with related data (Job, Company, User, Profile)
	var applications []models.Application
	if err := database.DB.
		Preload("User.Profile").    // Preload User and their Profile
		Preload("Job.Company").     // Preload Job and its associated Company
		Where("user_id = ?", userID).
		Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	// Transform data into the ApplicationResponse format
	var responses []models.ApplicationResponse
	for _, app := range applications {
		// Build the response for each application
		response := models.ApplicationResponse{
			ID:        app.ID,
			Status:    app.Status,
			CreatedAt: app.CreatedAt,
			UpdatedAt: app.UpdatedAt,
		}

		// Populate user details
		if app.User.ID != 0 {
			response.User.Name = app.User.Name
			response.User.Email = app.User.Email

			// Populate profile details
			response.Profile.Bio = app.User.Profile.Bio
			response.Profile.Resume = app.User.Profile.Resume
			response.Profile.Skills = app.User.Profile.Skills
		}

		// Populate company details
		if app.Job.ID != 0 && app.Job.Company.ID != 0 {
			response.Company.ID = app.Job.CompanyID
			response.Company.Name = app.Job.Company.Name
		}

		// Add the response to the list
		responses = append(responses, response)
	}

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}



func UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Decode and validate the payload
	var payload struct {
		ApplicationIDs []uint `json:"application_ids"` // List of application IDs
		Status         string `json:"status"`         // Desired status
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(payload.ApplicationIDs) == 0 {
		http.Error(w, "No application IDs provided", http.StatusBadRequest)
		return
	}

	// Validate the status value
	allowedStatuses := map[string]bool{
		"Pending":  true,
		"Accepted": true,
		"Rejected": true,
	}
	if !allowedStatuses[payload.Status] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Get the employer's company
	userID := uint(claims["user_id"].(float64))
	var company models.Company
	if err := database.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		http.Error(w, "Failed to find employer's company", http.StatusInternalServerError)
		return
	}

	// Fetch the applications to ensure they belong to the employer's jobs
	var applications []models.Application
	if err := database.DB.Joins("JOIN jobs ON applications.job_id = jobs.id").
		Where("applications.id IN ? AND jobs.company_id = ?", payload.ApplicationIDs, company.ID).
		Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	// Check if all applications were found
	if len(applications) != len(payload.ApplicationIDs) {
		http.Error(w, "Some applications do not belong to your jobs or do not exist", http.StatusForbidden)
		return
	}

	// Update the status of the applications in batch
	if err := database.DB.Model(&models.Application{}).
		Where("id IN ?", payload.ApplicationIDs).
		Update("status", payload.Status).Error; err != nil {
		http.Error(w, "Failed to update application statuses", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message"			: "Application statuses updated successfully",
		"updated_status"	: payload.Status,
		"application_ids"	: payload.ApplicationIDs,
	})
}

