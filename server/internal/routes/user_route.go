package routes

import (
	"server/internal/handlers"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, handler *handlers.UserHandler) {
	user := r.Group("/api/v1/user", middleware.AuthRequired(), middleware.RoleOnly("user"))

	user.GET("/profile", handler.GetUserProfile)
	user.PUT("/profile", handler.UpdateUserProfile)
	user.GET("/subscriptions", handler.GetMySubscription)
	user.PUT("/subscriptions", handler.UpdateUserSubscription)
	user.GET("/payments", handler.GetMyTransactionHistory)
	user.POST("/forms", handler.GetMyForms)
	user.POST("/forms/:id", handler.GetMyFormDetail)

}
