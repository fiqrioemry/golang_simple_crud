package middleware

import (
	"net/http"
	"server/internal/utils"
	"slices"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("accessToken")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized!! Token missing"})
			return
		}

		claims, err := utils.DecodeAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func RoleOnly(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := utils.MustGetRole(c)
		if slices.Contains(allowedRoles, role) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden: Access denied"})
	}
}
