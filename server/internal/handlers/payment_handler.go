package handlers

import (
	"net/http"
	"server/internal/dto"
	"server/internal/services"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service}
}

func (h *PaymentHandler) CreateNewPayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	userID := utils.MustGetUserID(c)

	res, err := h.service.CreatePayment(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create payment", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *PaymentHandler) HandlePaymentNotification(c *gin.Context) {
	var notif dto.MidtransNotificationRequest
	if !utils.BindAndValidateJSON(c, &notif) {
		return
	}

	if err := h.service.HandlePaymentNotification(notif); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process payment notification", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successfully processed"})
}

func (h *PaymentHandler) GetAllPaymentHistory(c *gin.Context) {
	query := c.Query("q")
	page := utils.StringToInt(c.DefaultQuery("page", "1"))
	limit := utils.StringToInt(c.DefaultQuery("limit", "10"))

	res, err := h.service.GetAllUserPayments(query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch payment history", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *PaymentHandler) GetPaymentDetail(c *gin.Context) {
	id := c.Param("id")
	payment, err := h.service.GetPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Payment not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
