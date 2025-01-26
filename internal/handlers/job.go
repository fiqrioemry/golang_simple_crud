package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateJob handles job creation for employers.
func CreateJob(w http.ResponseWriter, r *http.Request) {
	// Define a struct to handle the request payload
	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Location    string   `json:"location"`
		Type        string   `json:"type"`
		Skills      []string `json:"skills"`
		Experience  string   `json:"experience"`
	}

	// Decode the request payload into the struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.Location == "" || req.Type == "" || req.Experience == "" || len(req.Skills) == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate `Type` field
	validTypes := map[string]bool{"fulltime": true, "contract": true, "freelance": true, "internship": true}
	if !validTypes[req.Type] {
		http.Error(w, "Invalid job type. Must be one of [fulltime, contract, freelance, internship]", http.StatusBadRequest)
		return
	}

	// Validate `Experience` field
	validExperience := map[string]bool{"fresh graduate": true, "0-5 tahun": true, "5-10 tahun": true}
	if !validExperience[req.Experience] {
		http.Error(w, "Invalid experience level. Must be one of [fresh graduate, 0-5 tahun, 5-10 tahun]", http.StatusBadRequest)
		return
	}

	// Get the employer's company ID from the logged-in user
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))
	var company models.Company
	if err := database.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		http.Error(w, "Failed to find employer's company", http.StatusInternalServerError)
		return
	}

	// Create the Job model
	job := models.Job{
		CompanyID:   company.ID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Type:        req.Type,
		Skills:      req.Skills,
		Experience:  req.Experience,
	}

	// Save the job to the database
	if err := database.DB.Create(&job).Error; err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Job created successfully"})
}


// GetAllJobs handles getting all available jobs from company
// GetAllJobs handles getting all available jobs from companies, including the number of applications for each job.
func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	var companies []models.Company

	// Preload jobs and applications for each company
	if err := database.DB.Preload("Jobs.Applications").Find(&companies).Error; err != nil {
		http.Error(w, "Failed to retrieve jobs from companies", http.StatusInternalServerError)
		return
	}

	// Transform the response to include the number of applications for each job
	response := []map[string]interface{}{}
	for _, company := range companies {
		companyData := map[string]interface{}{
			"id":          company.ID,
			"name":        company.Name,
			"description": company.Description,
			"location":    company.Location,
			"created_at":  company.CreatedAt,
			"updated_at":  company.UpdatedAt,
			"jobs":        []map[string]interface{}{},
		}

		for _, job := range company.Jobs {
			jobData := map[string]interface{}{
				"id":           job.ID,
				"title":        job.Title,
				"description":  job.Description,
				"location":     job.Location,
				"type":         job.Type,
				"skills":       job.Skills,
				"experience":   job.Experience,
				"applications": len(job.Applications), // Count applications
				"created_at":   job.CreatedAt,
				"updated_at":   job.UpdatedAt,
			}
			companyData["jobs"] = append(companyData["jobs"].([]map[string]interface{}), jobData)
		}

		response = append(response, companyData)
	}

	// Respond with the transformed data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetJobByID fetches a specific job and includes only the company name.
func GetJobByID(w http.ResponseWriter, r *http.Request) {
	jobID := mux.Vars(r)["id"]

	var job models.Job
	if err := database.DB.Preload("Company").First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Transform the data
	jobResponse := models.JobResponse{
		ID:          job.ID,
		Title:       job.Title,
		Description: job.Description,
		Location:    job.Location,
		Type:        job.Type,
		Skills:      job.Skills,
		Experience:  job.Experience,
		CompanyName: job.Company.Name,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobResponse)
}



// UpdateJob handles updating a job posting (employer only).
func UpdateJob(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Get the job ID from the request URL
	jobID := mux.Vars(r)["id"]

	// Find the job in the database
	var job models.Job
	if err := database.DB.First(&job, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	// Verify that the job belongs to the employer's company
	userID := uint(claims["user_id"].(float64))
	var company models.Company
	if err := database.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		http.Error(w, "Failed to find employer's company", http.StatusInternalServerError)
		return
	}

	if job.CompanyID != company.ID {
		http.Error(w, "Unauthorized: Cannot update job of another employer", http.StatusForbidden)
		return
	}

	// Define a struct to handle the request payload
	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Location    string   `json:"location"`
		Type        string   `json:"type"`
		Skills      []string `json:"skills"`
		Experience  string   `json:"experience"`
	}

	// Decode the request payload into the struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields (if provided)
	if req.Title == "" || req.Description == "" || req.Location == "" || req.Type == "" || req.Experience == "" || len(req.Skills) == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate `Type` field
	validTypes := map[string]bool{"fulltime": true, "contract": true, "freelance": true, "internship": true}
	if !validTypes[req.Type] {
		http.Error(w, "Invalid job type. Must be one of [fulltime, contract, freelance, internship]", http.StatusBadRequest)
		return
	}

	// Validate `Experience` field
	validExperience := map[string]bool{"fresh graduate": true, "0-5 tahun": true, "5-10 tahun": true}
	if !validExperience[req.Experience] {
		http.Error(w, "Invalid experience level. Must be one of [fresh graduate, 0-5 tahun, 5-10 tahun]", http.StatusBadRequest)
		return
	}

	// Update the job fields
	job.Title = req.Title
	job.Description = req.Description
	job.Location = req.Location
	job.Type = req.Type
	job.Skills = req.Skills
	job.Experience = req.Experience

	// Save the updated job to the database
	if err := database.DB.Save(&job).Error; err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Job updated successfully"})
}



// DeleteJob handles deleting a job posting (employer only).
func DeleteJob(w http.ResponseWriter, r *http.Request) {
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

	// Only allow the employer who created the job to delete it
	userID := claims["user_id"].(float64)
	var company models.Company
	if err := database.DB.Where("user_id = ?", uint(userID)).First(&company).Error; err != nil {
		http.Error(w, "Failed to find employer's company", http.StatusInternalServerError)
		return
	}

	if job.CompanyID != company.ID {
		http.Error(w, "Unauthorized: Cannot delete job of another employer", http.StatusForbidden)
		return
	}

	// Delete the job
	if err := database.DB.Delete(&job).Error; err != nil {
		http.Error(w, "Failed to delete job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Job deleted successfully"})
}
