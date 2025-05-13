package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	service services.FormService
}

func NewFormHandler(service services.FormService) *FormHandler {
	return &FormHandler{service}
}

func (h *FormHandler) CreateNewForm(c *gin.Context) {
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
}

func (h *FormHandler) GetAllForms(c *gin.Context) {
}

func (h *FormHandler) GetFormDetail(c *gin.Context) {
}

func (h *FormHandler) DeleteForm(c *gin.Context) {
}
