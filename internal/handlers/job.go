package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateJob(w http.ResponseWriter, r *http.Request) {
	// Deklarasi struct untuk request payload
	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Location    string   `json:"location"`
		Type        string   `json:"type"`
		Skills      []string `json:"skills"`
		Experience  string   `json:"experience"`
	}

	// Decode JSON dari request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validasi field kosong
	if req.Title == "" || req.Description == "" || req.Location == "" || req.Type == "" || req.Experience == "" || len(req.Skills) == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validasi job type
	validTypes := map[string]bool{"fulltime": true, "contract": true, "freelance": true, "internship": true}
	if !validTypes[req.Type] {
		http.Error(w, "Invalid job type. Must be one of [fulltime, contract, freelance, internship]", http.StatusBadRequest)
		return
	}

	// Validasi experience level
	validExperience := map[string]bool{"fresh graduate": true, "0-5 tahun": true, "5-10 tahun": true}
	if !validExperience[req.Experience] {
		http.Error(w, "Invalid experience level. Must be one of [fresh graduate, 0-5 tahun, 5-10 tahun]", http.StatusBadRequest)
		return
	}

	// Ambil klaim dari middleware JWT
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Validasi dan konversi company_id dari claims
	companyIDFloat, ok := claims["company_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims: missing company_id", http.StatusUnauthorized)
		return
	}
	companyID := uint(companyIDFloat)

	// Buat job baru
	job := models.Job{
		CompanyID:   companyID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Type:        req.Type,
		Skills:      req.Skills,
		Experience:  req.Experience,
	}

	if err := database.DB.Create(&job).Error; err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Job created successfully",
		"job_id":  job.ID,
	})
}

// func CreateJob(w http.ResponseWriter, r *http.Request) {

// 	var req struct {
// 		Title       string   `json:"title"`
// 		Description string   `json:"description"`
// 		Location    string   `json:"location"`
// 		Type        string   `json:"type"`
// 		Skills      []string `json:"skills"`
// 		Experience  string   `json:"experience"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Title == "" || req.Description == "" || req.Location == "" || req.Type == "" || req.Experience == "" || len(req.Skills) == 0 {
// 		http.Error(w, "All fields are required", http.StatusBadRequest)
// 		return
// 	}

// 	validTypes := map[string]bool{"fulltime": true, "contract": true, "freelance": true, "internship": true}
// 	if !validTypes[req.Type] {
// 		http.Error(w, "Invalid job type. Must be one of [fulltime, contract, freelance, internship]", http.StatusBadRequest)
// 		return
// 	}

// 	validExperience := map[string]bool{"fresh graduate": true, "0-5 tahun": true, "5-10 tahun": true}
// 	if !validExperience[req.Experience] {
// 		http.Error(w, "Invalid experience level. Must be one of [fresh graduate, 0-5 tahun, 5-10 tahun]", http.StatusBadRequest)
// 		return
// 	}

// 	claims, err := middleware.GetUserFromContext(r)
// 	if err != nil || claims["role"] != "employer" {
// 		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
// 		return
// 	}

// 	companyID := uint(claims["company_id"].(float64))

// 	job := models.Job{
// 		CompanyID:   companyID,
// 		Title:       req.Title,
// 		Description: req.Description,
// 		Location:    req.Location,
// 		Type:        req.Type,
// 		Skills:      req.Skills,
// 		Experience:  req.Experience,
// 	}

// 	if err := database.DB.Create(&job).Error; err != nil {
// 		http.Error(w, "Failed to create job", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Job created successfully"})
// }

func UpdateJob(w http.ResponseWriter, r *http.Request) {
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

	companyID := uint(claims["company_id"].(float64))

	if job.CompanyID != companyID {
		http.Error(w, "Unauthorized: Cannot update job of another employer", http.StatusForbidden)
		return
	}

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Location    string   `json:"location"`
		Type        string   `json:"type"`
		Skills      []string `json:"skills"`
		Experience  string   `json:"experience"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.Location == "" || req.Type == "" || req.Experience == "" || len(req.Skills) == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

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

	job.Title = req.Title
	job.Description = req.Description
	job.Location = req.Location
	job.Type = req.Type
	job.Skills = req.Skills
	job.Experience = req.Experience

	if err := database.DB.Save(&job).Error; err != nil {
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Job updated successfully"})
}

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

func GetAllJobs(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Job

	if err := database.DB.Preload("Applications").Preload("Company").Find(&jobs).Error; err != nil {
		http.Error(w, "Failed to retrieve jobs", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}
	for _, job := range jobs {
		jobData := map[string]interface{}{
			"id":                 job.ID,
			"title":              job.Title,
			"description":        job.Description,
			"location":           job.Location,
			"type":               job.Type,
			"skills":             job.Skills,
			"experience":         job.Experience,
			"created_at":         job.CreatedAt,
			"updated_at":         job.UpdatedAt,
			"company_id":         job.Company.ID,
			"company_name":       job.Company.Name,
			"total_applications": len(job.Applications),
		}
		response = append(response, jobData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetJobByID(w http.ResponseWriter, r *http.Request) {
	jobID := mux.Vars(r)["id"]

	var jobs []models.Job
	if err := database.DB.Preload("Applications").Preload("Company").First(&jobs, jobID).Error; err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	var response []map[string]interface{}
	for _, job := range jobs {
		jobData := map[string]interface{}{
			"id":          job.ID,
			"title":       job.Title,
			"description": job.Description,
			"location":    job.Location,
			"type":        job.Type,
			"skills":      job.Skills,
			"experience":  job.Experience,
			"created_at":  job.CreatedAt,
			"updated_at":  job.UpdatedAt,
			"company": map[string]interface{}{
				"id":          job.Company.ID,
				"name":        job.Company.Name,
				"description": job.Company.Description,
				"Location":    job.Company.Location,
			},
			"total_applications": len(job.Applications),
		}
		response = append(response, jobData)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllEmployerPostedJobs(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	var jobs []models.Job

	companyID := uint(claims["company_id"].(float64))
	if err := database.DB.Preload("Applications").Where("company_id = ?", companyID).First(&jobs).Error; err != nil {
		http.Error(w, "Failed to retrieve jobs", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}
	for _, job := range jobs {
		applicationData := map[string]interface{}{
			"id":                 job.ID,
			"title":              job.Title,
			"type":               job.Type,
			"Location":           job.Location,
			"total_applications": len(job.Applications),
			"created_at":         job.CreatedAt,
			"updated_at":         job.UpdatedAt,
		}
		response = append(response, applicationData)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
