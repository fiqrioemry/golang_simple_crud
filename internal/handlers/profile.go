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
	if err := database.DB.Preload("Profile").First(&user, userID).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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
		Company   string `gorm:"size:100"`
		Title     string `gorm:"size:100"`
		StartDate time.Time
		EndDate   *time.Time
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

	var experience models.Experience

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Experience added successfully",
		"experience": experience,
	})
}

func UpdateUserSeekerExperience(w http.ResponseWriter, r *http.Request) {
	// Ambil klaim user dari context untuk memverifikasi role dan user_id
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	// Decode request payload
	var payload struct {
		Company   string     `json:"company"`
		Title     string     `json:"title"`
		StartDate time.Time  `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validasi input
	if payload.Company == "" || payload.Title == "" || payload.StartDate.IsZero() {
		http.Error(w, "Company, Title, and Start Date are required fields", http.StatusBadRequest)
		return
	}
	if payload.EndDate != nil && payload.EndDate.Before(payload.StartDate) {
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

	experience.Company = payload.Company
	experience.Title = payload.Title
	experience.StartDate = payload.StartDate
	experience.EndDate = payload.EndDate

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
