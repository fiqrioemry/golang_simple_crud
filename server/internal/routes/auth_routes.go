package routes

import (
	"server/internal/handlers"

	"github.com/fiqrioemry/microservice-ecommerce/server/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, handler *handlers.AuthHandler) {
	auth := r.Group("/api/v1/auth")

	auth.POST("/send-otp", handler.SendOTP)
	auth.POST("/verify-otp", handler.VerifyOTP)
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	protected := auth.Use(middleware.AuthRequired())
	protected.POST("/me", handler.AuthMe)
}
