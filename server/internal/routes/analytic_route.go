package routes

import (
	"server/internal/handlers"

	"github.com/fiqrioemry/microservice-ecommerce/server/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func AnalyticsRoutes(r *gin.Engine, handler *handlers.AnalyticsHandler) {
	analytics := r.Group("/api/v1/forms", middleware.AuthRequired(), middleware.RoleOnly("user", "admin"))

	analytics.GET("/:id/analytics", handler.GetFormAnalytics)
	analytics.GET("/:id/analytics/summary", handler.GetFormAnalyticSummary)
}
