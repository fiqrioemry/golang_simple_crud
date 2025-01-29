package handlers

import (
	"encoding/json"
	"golang_project/internal/auth"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register : Seeker
func SeekerRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusUnprocessableEntity)
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {

		user := models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			Role:     "seeker",
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Create the seeker profile
		seeker := models.Seeker{
			UserID: user.ID,
			Name:   req.Name,
		}

		if err := tx.Create(&seeker).Error; err != nil {
			return err // Transaction will automatically rollback on error
		}

		return nil
	}); err != nil {
		http.Error(w, "Failed to register new seeker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func EmployerRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Location string `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.Location = strings.TrimSpace(req.Location)

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Location == "" {
		http.Error(w, "All fields are required", http.StatusUnprocessableEntity)
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {

		user := models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			Role:     "employer",
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		employer := models.Employer{
			UserID:   user.ID,
			Name:     req.Name,
			Location: req.Location,
		}

		if err := tx.Create(&employer).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		http.Error(w, "Failed to register new employer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employer registered successfully"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	credentials.Email = strings.TrimSpace(credentials.Email)
	credentials.Password = strings.TrimSpace(credentials.Password)

	if credentials.Email == "" || credentials.Password == "" {
		http.Error(w, "Email and password are required", http.StatusUnprocessableEntity)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, string(user.Role), 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateToken(user.ID, string(user.Role), 7*24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Login successful",
		"access_token": accessToken,
	})
}

func AuthMe(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	userID := uint(userIDFloat)

	userRole, ok := claims["role"].(string)
	if !ok {
		http.Error(w, "Invalid role claims", http.StatusUnauthorized)
		return
	}

	var response map[string]interface{}

	if userRole == "seeker" {
		var seeker models.Seeker
		if err := database.DB.Preload("User").Where("user_id = ?", userID).First(&seeker).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		response = map[string]interface{}{
			"user_id": seeker.UserID,
			"email":   seeker.User.Email,
			"name":    seeker.Name,
			"role":    userRole,
		}

	} else if userRole == "employer" {
		var employer models.Employer
		if err := database.DB.Preload("User").Where("user_id = ?", userID).First(&employer).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		response = map[string]interface{}{
			"user_id":      employer.UserID,
			"email":        employer.User.Email,
			"role":         userRole,
			"company_id":   employer.ID,
			"company_name": employer.Name,
		}

	} else {
		http.Error(w, "Invalid user role", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateToken(uint(claims["user_id"].(float64)), claims["role"].(string), 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}
