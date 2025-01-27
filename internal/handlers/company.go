package handlers

import (
	"encoding/json"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"
)

func GetUserEmployerProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := database.DB.Preload("Company").First(&user, userID).Error; err != nil {
		http.Error(w, "User Company profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func EditUserEmployerProfile(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil || claims["role"] != "employer" {
		http.Error(w, "Unauthorized: Employer only", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Location    string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := database.DB.Preload("Company").First(&user, userID).Error; err != nil {
		http.Error(w, "User Company profile not found", http.StatusNotFound)
		return
	}

	user.Company.Name = req.Name
	user.Company.Description = req.Description
	user.Company.Location = req.Location

	if err := database.DB.Save(&user.Company).Error; err != nil {
		http.Error(w, "Failed to update company profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Company profile updated successfully",
		"payload": user.Company,
	})
}
