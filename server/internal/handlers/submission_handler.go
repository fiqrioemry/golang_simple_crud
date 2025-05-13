package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	service services.SubmissionService
}

func NewSubmissionHandler(service services.SubmissionService) *SubmissionHandler {
	return &SubmissionHandler{service}
}

func (h *SubmissionHandler) SendFormSubmission(c *gin.Context)   {}
func (h *SubmissionHandler) GetFormSubmissions(c *gin.Context)   {}
func (h *SubmissionHandler) GetSubmissionsResult(c *gin.Context) {}
