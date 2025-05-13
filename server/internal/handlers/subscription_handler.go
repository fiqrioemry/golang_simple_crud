package handlers

import (
	"net/http"
	"server/internal/dto"
	"server/internal/services"

	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service services.SubscriptionService
}

func NewSubscriptionHandler(service services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service}
}

// 1. Create Tier
func (h *SubscriptionHandler) CreateSubscriptionTier(c *gin.Context) {
	var req dto.CreateTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	if err := h.service.CreateTier(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create tier", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tier created successfully"})
}

// 2. Update Tier
func (h *SubscriptionHandler) UpdateSubscriptionTier(c *gin.Context) {
	var req dto.UpdateTierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	if err := h.service.UpdateTier(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update tier", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tier updated successfully"})
}

// 3. Delete Tier
func (h *SubscriptionHandler) DeleteSubscriptionTier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tier ID"})
		return
	}

	if err := h.service.DeleteTier(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete tier", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tier deleted successfully"})
}

// 4. Reset User Token (via cron job)
func (h *SubscriptionHandler) ResetUserSubscription(c *gin.Context) {
	if err := h.service.ResetAllUserTokens(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reset user tokens", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All active user tokens reset successfully"})
}

// 5. Get All Users + Subscription
func (h *SubscriptionHandler) GetAllUsersWithSubscriptions(c *gin.Context) {
	data, err := h.service.GetAllSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch subscriptions", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 6. Get One User Subscription Detail
func (h *SubscriptionHandler) GetUserDetailSubscriptions(c *gin.Context) {
	userID := c.Param("id")

	data, err := h.service.GetSubscriptionByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User subscription not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
