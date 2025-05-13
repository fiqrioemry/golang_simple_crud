package routes

import (
	"server/internal/handlers"

	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func QueueRoutes(r *gin.Engine, handler *handlers.QueueHandler) {
	queue := r.Group("/api/v1/queue", middleware.AuthRequired(), middleware.RoleOnly("user", "admin"))

	queue.GET("", handler.GetAllQueue)
	queue.POST("/:responseId/execute", handler.ExecuteQueue)
	queue.POST("/:responseId/complete", handler.CompleteQueue)
}
