package handlers

import (
	"server/internal/dto"
	"server/internal/services"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	res, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	fullname := c.PostForm("fullname")
	if len(fullname) < 3 {
		c.JSON(400, gin.H{"message": "Fullname must be at least 3 characters"})
		return
	}

	file, _ := c.FormFile("avatar")

	url := ""
	var err error
	if file != nil {
		url, err = utils.UploadImageWithValidation(file)
		if err != nil {
			c.JSON(400, gin.H{"message": "Invalid avatar file", "error": err.Error()})
			return
		}
	}

	err = h.service.UpdateProfile(userID, &dto.UpdateProfileRequest{
		Fullname: fullname,
		Avatar:   url,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to update profile", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Profile updated successfully", "avatar": url})
}

func (h *UserHandler) GetMySubscription(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	data, err := h.service.GetMySubscription(userID)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *UserHandler) UpdateUserSubscription(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	used := utils.GetQueryInt(c, "used", 1) // default: 1
	if err := h.service.UpdateUserSubscriptionToken(userID, used); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Token usage updated"})
}

func (h *UserHandler) GetMyTransactionHistory(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	data, err := h.service.GetTransactionHistory(userID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *UserHandler) GetMyForms(c *gin.Context) {
	userID := utils.MustGetUserID(c)
	data, err := h.service.GetMyForms(userID)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *UserHandler) GetMyFormDetail(c *gin.Context) {
	formID := c.Param("id")
	data, err := h.service.GetMyFormDetail(formID)
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, data)
}
