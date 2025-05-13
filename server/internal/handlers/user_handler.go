package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.AuthService
}

func NewUserHandler(service services.AuthService) *UserHandler {
	return &UserHandler{service}
}

// get api/v1/user/profile
// mengambil informasi data user
func (h *UserHandler) GetUserProfile(c *gin.Context) {
}

// put api/v1/user/profile
// mengedit informasi data user seperti fullname, avatar
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
}

// get api/v1/user/subscriptions
// mengambil informasi dari user subscriptions
func (h *UserHandler) GetMySubscription(c *gin.Context) {
}

// post api/v1/user/subscriptions
// user menambahkan subscription ketika berhasil melakukan pembelian (mungkin bisa di transaksi saja)
func (h *UserHandler) CreateNewUserSubscription(c *gin.Context) {
}

//	============

// put api/v1/user/subscriptions
// mengurangi token pengguna setiap kali fitur AI digunakan
func (h *UserHandler) UpdateUserSubscription(c *gin.Context) {
}

// men-set ulang limit token harian untuk cron-job
func (h *SubscriptionHandler) ResetUserSubscription(c *gin.Context) {
}
