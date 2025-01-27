package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetUserSeekerProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := database.DB.Preload("Profile.Experience").Preload("Applications").First(&user, userID).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	var experienceData []map[string]interface{}
	for _, experience := range user.Profile.Experience {
		experienceData = append(experienceData, map[string]interface{}{
			"id":         experience.ID,
			"company":    experience.Company,
			"title":      experience.Title,
			"start_date": experience.StartDate,
			"end_date":   experience.EndDate,
		})
	}

	response := map[string]interface{}{
		"user_id":      user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"bio":          user.Profile.Bio,
		"resume":       user.Profile.Resume,
		"skills":       user.Profile.Skills,
		"experience":   experienceData,
		"applications": len(user.Applications),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUserSeekerProfile(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name   string   `json:"name"`
		Bio    string   `json:"bio"`
		Resume string   `json:"resume"`
		Skills []string `json:"skills"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok || uint(userID) == 0 {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := database.DB.Preload("Profile").First(&user, uint(userID)).Error; err != nil {
		http.Error(w, "User or profile not found", http.StatusNotFound)
		return
	}
	user.Name = req.Name
	user.Profile.Bio = req.Bio
	user.Profile.Resume = req.Resume
	user.Profile.Skills = req.Skills

	if err := database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if err := tx.Save(&user.Profile).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
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

	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		experience := models.Experience{
			ProfileID: profile.ID,
			Company:   req.Company,
			Title:     req.Title,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		}

		if err := tx.Create(&experience).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		http.Error(w, "Failed to Add new Experience", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Experience added successfully",
	})
}

func UpdateUserSeekerExperience(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

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

	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := database.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	experienceID := mux.Vars(r)["id"]

	var experience models.Experience
	if err := database.DB.Where("id = ? AND profile_id = ?", experienceID, user.Profile.ID).First(&experience).Error; err != nil {
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
		"message":    "Experience updated successfully",
		"experience": experience,
	})
}
