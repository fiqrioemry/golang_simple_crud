package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	service services.QueueService
}

func NewQueueHandler(service services.QueueService) *QueueHandler {
	return &QueueHandler{service}
}

func (h *QueueHandler) GetAllQueue(c *gin.Context)   {}
func (h *QueueHandler) ExecuteQueue(c *gin.Context)  {}
func (h *QueueHandler) CompleteQueue(c *gin.Context) {}
