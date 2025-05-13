package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service services.AdminService
}

func NewAdminHandler(service services.AdminService) *AdminHandler {
	return &AdminHandler{service}
}

// admin membuat tier subscription baru yang dapat dibeli
func (h *AdminHandler) CreateSubscriptionTier(c *gin.Context) {
}

// admin mengupdate parameter tier subscription yang dapat dibeli
func (h *AdminHandler) UpdateSubscriptionTier(c *gin.Context) {
}

// admin menghapus tier subscription yang dapat dibeli
func (h *AdminHandler) DeleteSubscriptionTier(c *gin.Context) {
}

// men-set ulang limit token harian untuk cron-job
func (h *AdminHandler) ResetUserSubscription(c *gin.Context) {
}

// admin melihat daftar user lengkap dengan informasi subscriptionsnya
func (h *AdminHandler) GetAllUsersWithSubscriptions(c *gin.Context) {
}

// admin melihat detail dari informasi subscription user
func (h *AdminHandler) GetUserDetailSubscriptions(c *gin.Context) {
}

func (h *AdminHandler) GetAllTransactionsHistory(c *gin.Context) {
}
