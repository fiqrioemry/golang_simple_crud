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

// Register : Seeker
func SeekerRegister(w http.ResponseWriter, r *http.Request) {

	var user models.User

	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check for required fields (e.g., Name, Email, and Password)
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All field are required", http.StatusBadRequest)
		return 
	}
	

	// Check if the email is already registered
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password 	= string(hashedPassword)
	user.Role = models.Seeker

	// Use a database transaction
	tx := database.DB.Begin() // Start transaction
	if tx.Error != nil {
		http.Error(w, "Failed to start database transaction", http.StatusInternalServerError)
		return
	}

	// Save user to the database
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback() // Rollback if error
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	profile := models.Profile{
		UserID : user.ID,
	}

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
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
	var user models.User
	var company models.Company

	// Decode the request body into the user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check required field
	if user.Name == "" || user.Email == "" || user.Password == "" || company.Name ==  "" || company.Description == "" || company.Location == "" {
		http.Error(w, "All field are required", http.StatusBadRequest)
		return
	}

	// Check if the email is already registered
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email is already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password 	= string(hashedPassword)
	user.Role = models.Employer

	// Use a database transaction
	tx := database.DB.Begin() // Start transaction
	if tx.Error != nil {
		http.Error(w, "Failed to start database transaction", http.StatusInternalServerError)
		return
	}

	// Save user to the database
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback() // Rollback the transaction if there is an error
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	company.UserID = user.ID

	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create company", http.StatusInternalServerError)
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


// Login handles user login and JWT generation.
func Login(w http.ResponseWriter, r *http.Request) {
	// Validasi metode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Dekode body request ke dalam struct credentials
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Cari pengguna berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or user not found", http.StatusUnauthorized)
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate akses token (JWT)
	accessToken, err := auth.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	refreshToken, err := auth.GenerateRefreshToken(user.ID, string(user.Role))
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Set refresh token sebagai cookie (opsional)
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour), // Refresh token berlaku 7 hari
		HttpOnly: true,                               // Hanya dapat diakses oleh server
		Secure:   false,                              // Set true jika menggunakan HTTPS
		Path:     "/",
	})

	// Kirimkan respons dengan akses token dan refresh token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Login is successful",
		"accessToken":  accessToken,
	})
}


func AuthMe(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	jsonResponse := map[string]interface{}{
		"user_id": claims["user_id"],
	}

	var user models.User
	if err := database.DB.Where("ID = ?", jsonResponse["user_id"]).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or user not found", http.StatusUnauthorized)
		return
	}

	payload := map[string]interface{
		"user_id" 	: user.ID,
		"email"		: user.Email,
		"name"		: user.Name
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"payload":  user,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Login is successful",
		"accessToken":  accessToken,
	})
}


