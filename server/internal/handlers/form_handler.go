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

func (h *FormHandler) CreateNewForm(c *gin.Context)      {}
func (h *FormHandler) GetAllForms(c *gin.Context)        {}
func (h *FormHandler) GetFormDetail(c *gin.Context)      {}
func (h *FormHandler) GetFormSettings(c *gin.Context)    {}
func (h *FormHandler) UpdateFormSettings(c *gin.Context) {}
func (h *FormHandler) AddFormSection(c *gin.Context)     {}
func (h *FormHandler) GetFormSections(c *gin.Context)    {}
func (h *FormHandler) UpdateFormSections(c *gin.Context) {}
func (h *FormHandler) DeleteFormSections(c *gin.Context) {}
func (h *FormHandler) GetFormQuestion(c *gin.Context)    {}
func (h *FormHandler) AddFormQuestion(c *gin.Context)    {}
func (h *FormHandler) UpdateQuestion(c *gin.Context)     {}
func (h *FormHandler) DeleteQuestion(c *gin.Context)     {}
