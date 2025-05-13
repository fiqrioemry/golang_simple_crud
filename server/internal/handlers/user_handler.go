package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) GetUserProfile(c *gin.Context)            {}
func (h *UserHandler) UpdateUserProfile(c *gin.Context)         {}
func (h *UserHandler) GetMySubscription(c *gin.Context)         {}
func (h *UserHandler) CreateNewUserSubscription(c *gin.Context) {}
func (h *UserHandler) UpdateUserSubscription(c *gin.Context)    {}
func (h *UserHandler) GetMyTransactionHistory(c *gin.Context)   {}
func (h *UserHandler) GetMyForms(c *gin.Context)                {}
func (h *UserHandler) GetMyFormDetail(c *gin.Context)           {}
