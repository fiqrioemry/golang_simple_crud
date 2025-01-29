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

	userID := uint(claims["user_id"].(float64))

	jobID := mux.Vars(r)["id"]

	var job models.Job
	if err := database.DB.First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	var existingApplication models.Application
	if err := database.DB.Where("job_id = ? AND user_id = ?", job.ID, userID).First(&existingApplication).Error; err == nil {
		http.Error(w, "You have already applied for this job", http.StatusConflict)
		return
	}

	application := models.Application{
		JobID:  job.ID,
		UserID: userID,
		Status: "Pending",
	}

	if err := database.DB.Create(&application).Error; err != nil {
		http.Error(w, "Failed to apply to job", http.StatusInternalServerError)
		return
	}

	// Kirimkan respons sukses
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
	if err := database.DB.Preload("Seeker.User").Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
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
			"name":       application.Seeker.Name,
			"email":      application.Seeker.User.Email,
			"status":     application.Status,
			"created_at": application.CreatedAt,
			"updated_at": application.UpdatedAt,
		}
		response = append(response, applicationData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetSeekerJobApplications(w http.ResponseWriter, r *http.Request) {

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
	if err := database.DB.Preload("Job.Applications").Preload("Job.Employer").Where("user_id = ?", uint(userID)).Find(&applications).Error; err != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
	}
	if len(applications) == 0 {
		http.Error(w, "User has no applications", http.StatusNotFound)
		return
	}

	var response []map[string]interface{}
	for _, applications := range applications {
		applicationData := map[string]interface{}{
			"application_id":     applications.ID,
			"job_id":             applications.JobID,
			"company_name":       applications.Job.Employer.Name,
			"job_title":          applications.Job.Title,
			"job_location":       applications.Job.Location,
			"total_applications": len(applications.Job.Applications),
			"status":             applications.Status,
			"created_at":         applications.CreatedAt,
		}
		response = append(response, applicationData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {

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

	allowedStatuses := map[string]bool{
		"Pending":  true,
		"Accepted": true,
		"Rejected": true,
	}
	if !allowedStatuses[payload.Status] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	jobID := mux.Vars(r)["id"]

	var job models.Job
	if err := database.DB.First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found or no longer exist", http.StatusNotFound)
		return
	}

	userID := uint(claims["user_id"].(float64))

	var employer models.Employer
	if err := database.DB.Where("user_id = ?", userID).First(&employer).Error; err != nil {
		http.Error(w, "Employer Profile not found", http.StatusInternalServerError)
		return
	}

	if job.EmployerID != employer.ID {
		http.Error(w, "Unauthorized: Cannot update job of another employer", http.StatusForbidden)
		return
	}

	if err := database.DB.Model(&models.Application{}).
		Where("id IN ?", payload.ApplicationIDs).
		Update("status", payload.Status).Error; err != nil {
		http.Error(w, "Failed to update application statuses", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Application statuses updated successfully",
	})
}
