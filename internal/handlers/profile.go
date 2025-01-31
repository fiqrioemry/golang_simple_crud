package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"golang_project/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type EmployerResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	Picture     string    `json:"picture"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}

func GetEmployerCompanyProfile(w http.ResponseWriter, r *http.Request) {
	employerID := mux.Vars(r)["id"]

	var employer models.Employer
	if err := database.DB.Preload("Jobs").First(&employer, employerID).Error; err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	var employerJobs []map[string]interface{}
	for _, job := range employer.Jobs {
		employerJobs = append(employerJobs, map[string]interface{}{
			"id":                 job.ID,
			"title":              job.Title,
			"type":               job.Type,
			"Location":           job.Location,
			"total_applications": len(job.Applications),
			"created_at":         job.CreatedAt,
			"updated_at":         job.UpdatedAt,
		})
	}

	response := map[string]interface{}{
		"company_id":   employer.ID,
		"company_name": employer.Name,
		"avatar":       employer.Avatar,
		"picture":      employer.Picture,
		"description":  employer.Description,
		"location":     employer.Location,
		"jobs":         employerJobs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetSeekerProfile(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	var seeker models.Seeker
	if err := database.DB.Preload("User").Preload("Experience").Preload("Applications").
		Where("user_id = ?", uint(userID)).Take(&seeker).Error; err != nil {
		http.Error(w, "Seeker profile not found", http.StatusNotFound)
		return
	}
	var seekerExperiences []map[string]interface{}

	for _, experience := range seeker.Experience {

		seekerExperience := map[string]interface{}{
			"experience_id": experience.ID,
			"company":       experience.Company,
			"title":         experience.Title,
			"startDate":     experience.StartDate,
			"endDate":       experience.EndDate,
		}

		seekerExperiences = append(seekerExperiences, seekerExperience)
	}

	response := map[string]interface{}{
		"user_id":      seeker.UserID,
		"name":         seeker.Name,
		"email":        seeker.User.Email,
		"bio":          seeker.Bio,
		"resume":       seeker.Resume,
		"skills":       seeker.Skills,
		"experience":   seekerExperiences,
		"applications": len(seeker.Applications),
		"role":         "seeker",
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateSeekerProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := uint(userIDFloat)

	err = r.ParseMultipartForm(10 << 20) // max 10mb
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var req struct {
		Name     string     `json:"name"`
		Bio      string     `json:"bio,omitempty"`
		Resume   string     `json:"resume,omitempty"`
		Skills   []string   `json:"skills,omitempty"`
		Gender   string     `json:"gender,omitempty"`
		Birthday *time.Time `json:"birthday,omitempty"`
	}

	if r.FormValue("name") == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(r.FormValue("name"))
	req.Bio = r.FormValue("bio")
	req.Resume = r.FormValue("resume")
	req.Gender = r.FormValue("gender")

	validGenders := map[string]bool{"male": true, "female": true}
	if req.Gender != "" && !validGenders[req.Gender] {
		http.Error(w, "Invalid gender. Must be one of [male, female]", http.StatusBadRequest)
		return
	}

	var avatarURL string
	file, fileHeader, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarURL, err = utils.UploadMediaToCloudinary(file, fileHeader)
		if err != nil {
			http.Error(w, "Failed to upload avatar to Cloudinary", http.StatusInternalServerError)
			return
		}
	}

	var seeker models.Seeker
	if err := database.DB.Where("user_id = ?", userID).Take(&seeker).Error; err != nil {
		http.Error(w, "Seeker Profile not found", http.StatusNotFound)
		return
	}

	seeker.Name = req.Name
	seeker.Bio = req.Bio
	seeker.Resume = req.Resume
	seeker.Skills = req.Skills
	seeker.Gender = req.Gender
	seeker.Birthday = req.Birthday

	if avatarURL != "" {
		seeker.Avatar = avatarURL
	}

	seeker.UpdatedAt = time.Now()

	if err := database.DB.Save(&seeker).Error; err != nil {
		http.Error(w, "Failed to update seeker profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Seeker profile updated successfully",
		"payload": seeker,
	})
}
func GetEmployerProfile(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := uint(userIDFloat)

	var employer models.Employer
	if err := database.DB.Preload("User").Preload("Jobs").
		Where("user_id = ?", userID).Take(&employer).Error; err != nil {
		http.Error(w, "Employer profile not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"user_id":      employer.UserID,
		"email":        employer.User.Email,
		"company_id":   employer.ID,
		"company_name": employer.Name,
		"location":     employer.Location,
		"jobs":         len(employer.Jobs),
	}

	json.NewEncoder(w).Encode(response)
}

func UpdateEmployerProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := uint(userIDFloat)

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
	}

	if r.FormValue("name") == "" {
		http.Error(w, "Company name is required", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(r.FormValue("name"))
	req.Description = r.FormValue("description")
	req.Location = r.FormValue("location")

	var employer models.Employer
	if err := database.DB.Where("user_id = ?", userID).Take(&employer).Error; err != nil {
		http.Error(w, "Employer Profile not found", http.StatusNotFound)
		return
	}

	var avatarURL string
	file, fileHeader, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		avatarURL, err = utils.UploadMediaToCloudinary(file, fileHeader)
		if err != nil {
			http.Error(w, "Failed to upload avatar to Cloudinary", http.StatusInternalServerError)
			return
		}
	}

	var pictureURL string
	file, fileHeader, err = r.FormFile("picture")
	if err == nil {
		defer file.Close()
		pictureURL, err = utils.UploadMediaToCloudinary(file, fileHeader)
		if err != nil {
			http.Error(w, "Failed to upload picture to Cloudinary", http.StatusInternalServerError)
			return
		}
	}

	employer.Name = req.Name
	employer.Location = req.Location
	employer.Description = req.Description
	if avatarURL != "" {
		employer.Avatar = avatarURL
	}
	if pictureURL != "" {
		employer.Picture = pictureURL
	}
	employer.UpdatedAt = time.Now()

	if err := database.DB.Save(&employer).Error; err != nil {
		http.Error(w, "Failed to update employer profile", http.StatusInternalServerError)
		return
	}

	response := EmployerResponse{
		ID:          employer.ID,
		CreatedAt:   employer.CreatedAt,
		UpdatedAt:   employer.UpdatedAt,
		UserID:      employer.UserID,
		Name:        employer.Name,
		Avatar:      employer.Avatar,
		Picture:     employer.Picture,
		Description: employer.Description,
		Location:    employer.Location,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Employer profile updated successfully",
		"payload": response,
	})
}

func AddUserSeekerExperience(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	var req struct {
		Company   string     `json:"company" gorm:"size:100"`
		Title     string     `json:"title" gorm:"size:100"`
		StartDate time.Time  `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Company == "" || req.Title == "" || req.StartDate.IsZero() {
		http.Error(w, "Company, Title, and Start Date are required fields", http.StatusBadRequest)
		return
	}

	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		http.Error(w, "End date cannot be before start date", http.StatusBadRequest)
		return
	}

	userID := uint(claims["user_id"].(float64))

	var seeker models.Seeker
	if err := database.DB.Where("user_id = ?", userID).Take(&seeker).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	experience := models.Experience{
		SeekerID:  seeker.ID,
		Company:   req.Company,
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := database.DB.Create(&experience).Error; err != nil {
		http.Error(w, "Failed to Add new Experience", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Experience added successfully",
		"payload": experience,
	})
}

func UpdateUserSeekerExperience(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	experienceID := mux.Vars(r)["id"]

	var req struct {
		Company   string     `json:"company"`
		Title     string     `json:"title"`
		StartDate time.Time  `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Company == "" || req.Title == "" || req.StartDate.IsZero() {
		http.Error(w, "Company, Title, and Start Date are required fields", http.StatusBadRequest)
		return
	}
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		http.Error(w, "End date cannot be before start date", http.StatusBadRequest)
		return
	}

	var experience models.Experience
	if err := database.DB.Take(&experience, experienceID).Error; err != nil {
		http.Error(w, "Experience not found", http.StatusNotFound)
		return
	}

	experience.Company = req.Company
	experience.Title = req.Title
	experience.StartDate = req.StartDate
	experience.EndDate = req.EndDate

	if err := database.DB.Save(&experience).Error; err != nil {
		http.Error(w, "Failed to update experience", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Experience updated successfully",
		"payload": experience,
	})
}
