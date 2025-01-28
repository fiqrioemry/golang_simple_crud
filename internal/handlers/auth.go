package handlers

import (
	"encoding/json"
	"golang_project/internal/auth"
	"golang_project/internal/database"
	"golang_project/internal/middleware"
	"golang_project/internal/models"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// func SeekerRegister(w http.ResponseWriter, r *http.Request) {
// 	// Define a struct to handle the request body
// 	var req struct {
// 		Name     string `json:"name"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	// Decode the request body into the struct
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	// Validate required fields
// 	if req.Name == "" || req.Email == "" || req.Password == "" {
// 		http.Error(w, "All fields are required", http.StatusBadRequest)
// 		return
// 	}

// 	// Check if the email is already registered
// 	var existingUser models.User
// 	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
// 		http.Error(w, "Email is already registered", http.StatusConflict)
// 		return
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	// Use a database transaction with automatic handling
// 	if err := database.DB.Transaction(func(tx *gorm.DB) error {
// 		// Create the user model
// 		user := models.User{
// 			Name:     req.Name,
// 			Email:    req.Email,
// 			Password: string(hashedPassword),
// 			Role:     "seeker", // Explicitly set the role
// 		}

// 		if err := tx.Create(&user).Error; err != nil {
// 			return err // Automatically triggers rollback
// 		}

// 		// Create the profile associated with the user
// 		profile := models.Profile{
// 			UserID: user.ID,
// 		}

// 		if err := tx.Create(&profile).Error; err != nil {
// 			return err // Automatically triggers rollback
// 		}

// 		// If everything succeeds, transaction will be committed
// 		return nil
// 	}); err != nil {
// 		http.Error(w, "Failed to register user", http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
// }

// Register : Seeker
func SeekerRegister(w http.ResponseWriter, r *http.Request) {
	// Define a struct to handle the request body
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if the email is already registered
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create the user model
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "seeker", // Explicitly set the role
	}

	// Use a database transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start database transaction", http.StatusInternalServerError)
		return
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	profile := models.Profile{
		UserID: user.ID,
	}

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Register : Employer
func EmployerRegister(w http.ResponseWriter, r *http.Request) {
	// Define a structure to handle the nested request body
	var req struct {
		User    models.User    `json:"user"`
		Company models.Company `json:"company"`
	}

	// Decode the request body into the struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.User.Name == "" || req.User.Email == "" || req.User.Password == "" ||
		req.Company.Name == "" || req.Company.Description == "" || req.Company.Location == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if the email is already registered
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.User.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	req.User.Password = string(hashedPassword)
	req.User.Role = "employer"

	// Use a database transaction
	tx := database.DB.Begin() // Start transaction
	if tx.Error != nil {
		http.Error(w, "Failed to start database transaction", http.StatusInternalServerError)
		return
	}

	// Save user to the database
	if err := tx.Create(&req.User).Error; err != nil {
		tx.Rollback() // Rollback the transaction if there is an error
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Associate the company with the user
	req.Company.UserID = req.User.ID

	// Save company to the database
	if err := tx.Create(&req.Company).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create company", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employer registered successfully"})
}

// Login handles user login and JWT generation.
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Preload("Company").Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate Access Token (valid for 1 day)
	accessToken, err := auth.GenerateToken(user.ID, user.Company.ID, string(user.Role), 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Company.ID, string(user.Role), 7*24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login is successfully", "access_token": accessToken})
}

// AuthMe handles the authenticated user info retrieval.
func AuthMe(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := database.DB.Preload("Company").Where("ID = ?", uint(userID)).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var companyID *uint
	var companyName *string
	if user.Company.ID != 0 {
		companyID = &user.Company.ID
		companyName = &user.Company.Name
	}

	response := map[string]interface{}{
		"user_id":      user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"role":         user.Role,
		"company_id":   companyID,
		"company_name": companyName,
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

	accessToken, err := auth.GenerateToken(uint(claims["user_id"].(float64)), uint(claims["company_id"].(float64)), claims["role"].(string), 24*time.Hour)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Respond with the new access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}
