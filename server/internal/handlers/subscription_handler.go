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

// admin membuat tier subscription baru yang dapat dibeli
func (h *SubscriptionHandler) CreateSubscriptionTier(c *gin.Context) {
}

// admin mengupdate parameter tier subscription yang dapat dibeli
func (h *SubscriptionHandler) UpdateSubscriptionTier(c *gin.Context) {
}

// admin menghapus tier subscription yang dapat dibeli
func (h *SubscriptionHandler) DeleteSubscriptionTier(c *gin.Context) {
}

// admin melihat daftar user yang subscribe
func (h *SubscriptionHandler) GetAllUserSubscription(c *gin.Context) {
}

// admin melihat detail dari informasi subscription user
func (h *SubscriptionHandler) GetSubscriptionInfo(c *gin.Context) {
}

// user menambahkan subscription ketika berhasil melakukan pembelian (mungkin bisa di transaksi saja)
func (h *SubscriptionHandler) CreateNewUserSubscription(c *gin.Context) {
}

// user menambahkan subscription ketika berhasil melakukan pembelian (mungkin bisa di transaksi saja)
func (h *SubscriptionHandler) UpdateUserSubscription(c *gin.Context) {
}
