package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"
)


func GetUserSeekerProfile(w http.ResponseWriter, r *http.Request) {
	// Check if the user is a seeker
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	// Extract user ID from JWT claims
	userID := uint(claims["user_id"].(float64))

	// Fetch the user's profile
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Respond with the profile data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateUserSeekerProfile(w http.ResponseWriter, r *http.Request) {
	// Check if the user is a seeker
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "seeker" {
		http.Error(w, "Unauthorized: Seeker only", http.StatusUnauthorized)
		return
	}

	// Extract user ID from JWT claims
	userID := uint(claims["user_id"].(float64))

	// Fetch the profile to ensure it exists
	var profile models.Profile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Decode the request payload
	var payload struct {
		Bio    string   `json:"bio"`
		Resume string   `json:"resume"`
		Skills []string `json:"skills"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the profile fields
	profile.Bio = payload.Bio
	profile.Resume = payload.Resume
	profile.Skills = payload.Skills

	// Save the updated profile to the database
	if err := database.DB.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profile updated successfully",
	})
}
