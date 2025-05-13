package handlers

import (
	"net/http"
	"server/internal/dto"
	"server/internal/services"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

// ===================== Analytics =====================
type AnalyticsHandler struct {
	service services.AnalyticsService
}

func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service}
}

func (h *AnalyticsHandler) GetFormAnalytics(c *gin.Context)       {}
func (h *AnalyticsHandler) GetFormAnalyticSummary(c *gin.Context) {}
