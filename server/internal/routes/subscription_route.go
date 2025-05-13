package routes

import (
	"server/internal/handlers"

	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SubscriptionRoutes(r *gin.Engine, handler *handlers.SubscriptionHandler) {
	admin := r.Group("/api/v1/subscriptions", middleware.AuthRequired(), middleware.RoleOnly("admin"))

	admin.GET("", handler.GetAllUsersWithSubscriptions)
	admin.GET("/:id", handler.GetUserDetailSubscriptions)
	admin.PUT("/reset", handler.ResetUserSubscription)

	admin.POST("", handler.CreateSubscriptionTier)
	admin.PUT("/:id", handler.UpdateSubscriptionTier)
	admin.DELETE("/:id", handler.DeleteSubscriptionTier)
}
