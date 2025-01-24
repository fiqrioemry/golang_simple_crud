package middleware

import (
	"context"
	"errors"
	"golang_project/internal/auth"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// UserContextKey is the key for user information in the request context.
type UserContextKey string

const (
	ContextUserKey UserContextKey = "user"
)

// JWTMiddleware validates the JWT token and adds user info to the request context.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract the token from the header
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := tokenParts[1]

		// Validate the token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add claims to the request context
		ctx := context.WithValue(r.Context(), ContextUserKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the JWT claims from the request context.
func GetUserFromContext(r *http.Request) (jwt.MapClaims, error) {
	claims, ok := r.Context().Value(ContextUserKey).(jwt.MapClaims)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return claims, nil
}
