package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"
)

func GetUserEmployerProfile(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Extract user ID from JWT claims
	userID := uint(claims["user_id"].(float64))

	// Fetch the company profile
	var company models.Company
	if err := database.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		http.Error(w, "Company profile not found", http.StatusNotFound)
		return
	}

	// Respond with the company data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}


func EditUserEmployerProfile(w http.ResponseWriter, r *http.Request) {
	// Check if the user is an employer
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	// Extract user ID from JWT claims
	userID := uint(claims["user_id"].(float64))

	// Fetch the company profile to ensure it exists
	var company models.Company
	if err := database.DB.Where("user_id = ?", userID).First(&company).Error; err != nil {
		http.Error(w, "Company profile not found", http.StatusNotFound)
		return
	}

	// Decode the request payload
	var payload struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the company fields
	company.Name = payload.Name
	company.Description = payload.Description
	company.Location = payload.Location

	// Save the updated company profile to the database
	if err := database.DB.Save(&company).Error; err != nil {
		http.Error(w, "Failed to update company profile", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Company profile updated successfully",
	})
}
