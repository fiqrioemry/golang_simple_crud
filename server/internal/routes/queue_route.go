package routes

import (
	"server/internal/handlers"

	"github.com/fiqrioemry/microservice-ecommerce/server/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func QueueRoutes(r *gin.Engine, handler *handlers.QueueHandler) {
	queue := r.Group("/api/v1/queue", middleware.AuthRequired(), middleware.RoleOnly("admin"))

	queue.GET("", handler.GetAllQueue)
	queue.POST("/:responseId/execute", handler.ExecuteQueue)
	queue.POST("/:responseId/complete", handler.CompleteQueue)
}
