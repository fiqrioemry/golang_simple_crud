package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service services.PaymentService
}

func NewPaymentHandler(service services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service}
}

func (h *PaymentHandler) CreateNewPayment(c *gin.Context)     {}
func (h *PaymentHandler) GetAllPaymentHistory(c *gin.Context) {}
func (h *PaymentHandler) GetPaymentDetail(c *gin.Context)     {}
