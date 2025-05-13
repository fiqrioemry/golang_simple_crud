package handlers

import (
	"net/http"
	"server/internal/dto"
	"server/internal/services"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Input Request", "error": err.Error()})
		return
	}
	tokens, err := h.service.UserRegister(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email Already Registered", "error": err.Error()})
		return
	}
	utils.SetAccessTokenCookie(c, tokens.AccessToken)
	utils.SetRefreshTokenCookie(c, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Register Successfully"})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Invalid Input Request", "error": err.Error()})
		return
	}

	tokens, err := h.service.UserLogin(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	utils.SetAccessTokenCookie(c, tokens.AccessToken)
	utils.SetAccessTokenCookie(c, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Login Successfully"})
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input request", "error": err.Error()})
		return
	}

	if err := h.service.SendOTP(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP send to email successfully"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input request", "error": err.Error()})
		return
	}

	if err := h.service.VerifyOTPCode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP Verified succesfully"})
}

func (h *AuthHandler) AuthMe(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	response, err := h.service.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized Access", "error": err.Error()})
	}

	response, err := h.service.RefreshUserToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	utils.SetAccessTokenCookie(c, response.AccessToken)

	utils.SetRefreshTokenCookie(c, response.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Refresh Successfully"})

}
