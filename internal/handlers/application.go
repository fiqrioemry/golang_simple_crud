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

func GetEmployerJobApplications(w http.ResponseWriter, r *http.Request) {
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

	var applications []models.Application
	if err := database.DB.Preload("User").Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}
	if len(applications) == 0 {
		http.Error(w, "No applications found for this job", http.StatusNotFound)
		return
	}

	var response []map[string]interface{}
	for _, application := range applications {
		applicationData := map[string]interface{}{
			"id":         application.ID,
			"user_id":    application.UserID,
			"name":       application.User.Name,
			"email":      application.User.Email,
			"status":     application.Status,
			"created_at": application.CreatedAt,
			"updated_at": application.UpdatedAt,
		}
		response = append(response, applicationData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetSeekerJobApplication(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Job seeker only", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	var applications []models.Application
	if err := database.DB.
		Preload("Job.Company").
		Where("user_id = ?", uint(userID)).
		Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	if len(applications) == 0 {
		http.Error(w, "User has no applications", http.StatusNotFound)
		return
	}

	var response []map[string]interface{}
	for _, application := range applications {
		response = append(response, map[string]interface{}{
			"application_id":    application.ID,
			"status":            application.Status,
			"company_id":        application.Job.CompanyID,
			"company_name":      application.Job.Company.Name,
			"job_id":            application.JobID,
			"title":             application.Job.Title,
			"location":          application.Job.Location,
			"total_application": len(application.Job.Applications),
			"created_at":        application.CreatedAt,
			"updated_at":        application.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Decode and validate the payload
	var payload struct {
		ApplicationIDs []uint `json:"application_ids"`
		Status         string `json:"status"`
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
		"message":         "Application statuses updated successfully",
		"updated_status":  payload.Status,
		"application_ids": payload.ApplicationIDs,
	})
}
