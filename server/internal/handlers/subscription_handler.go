package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service services.SubscriptionService
}

func NewSubscriptionHandler(service services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service}
}

func (h *SubscriptionHandler) CreateSubscriptionTier(c *gin.Context)       {}
func (h *SubscriptionHandler) UpdateSubscriptionTier(c *gin.Context)       {}
func (h *SubscriptionHandler) DeleteSubscriptionTier(c *gin.Context)       {}
func (h *SubscriptionHandler) ResetUserSubscription(c *gin.Context)        {}
func (h *SubscriptionHandler) GetAllUsersWithSubscriptions(c *gin.Context) {}
func (h *SubscriptionHandler) GetUserDetailSubscriptions(c *gin.Context)   {}
