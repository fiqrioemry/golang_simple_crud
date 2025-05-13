package handlers

import (
	"server/internal/dto"
	"server/internal/services"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	service services.FormService
}

func NewFormHandler(service services.FormService) *FormHandler {
	return &FormHandler{service}
}

func (h *FormHandler) CreateNewForm(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	var req dto.CreateFormRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.CreateForm(userID, &req); err != nil {
		c.JSON(500, gin.H{"message": "Failed to create form", "error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Form created successfully"})
}

func (h *FormHandler) GetAllForms(c *gin.Context) {
	userID := utils.MustGetUserID(c)

	forms, err := h.service.GetAllForms(userID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to fetch forms", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": forms})
}

func (h *FormHandler) GetFormDetail(c *gin.Context) {
	formID := c.Param("formId")

	data, err := h.service.GetFormDetail(formID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Form not found", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": data})
}

func (h *FormHandler) GetFormSettings(c *gin.Context) {
	formID := c.Param("formId")

	data, err := h.service.GetFormSettings(formID)
	if err != nil {
		c.JSON(404, gin.H{"message": "Form setting not found", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h *FormHandler) UpdateFormSettings(c *gin.Context) {
	formID := c.Param("formId")

	var req dto.UpdateFormSettingRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.service.UpdateFormSettings(formID, &req); err != nil {
		c.JSON(500, gin.H{"message": "Failed to update form setting", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Form setting updated successfully"})
}

func (h *FormHandler) AddFormSection(c *gin.Context) {
	var req dto.AddSectionRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}

	data, err := h.service.AddFormSection(&req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to add section", "error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Section added successfully", "data": data})
}

func (h *FormHandler) GetFormSections(c *gin.Context) {
	formID := c.Param("formId")
	data, err := h.service.GetFormSections(formID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to fetch sections", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h *FormHandler) UpdateFormSections(c *gin.Context) {
	var req dto.UpdateSectionRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}
	if err := h.service.UpdateFormSections(&req); err != nil {
		c.JSON(500, gin.H{"message": "Failed to update section", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Section updated successfully"})
}

func (h *FormHandler) DeleteFormSections(c *gin.Context) {
	sectionID := c.Param("sectionId")
	if err := h.service.DeleteFormSections(sectionID); err != nil {
		c.JSON(500, gin.H{"message": "Failed to delete section", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Section deleted successfully"})
}

func (h *FormHandler) GetFormQuestion(c *gin.Context) {
	formID := c.Param("formId")
	data, err := h.service.GetFormQuestion(formID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to fetch questions", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h *FormHandler) AddFormQuestion(c *gin.Context) {
	var req dto.AddQuestionRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}
	if err := h.service.AddFormQuestion(&req); err != nil {
		c.JSON(500, gin.H{"message": "Failed to add question", "error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Question added successfully"})
}

func (h *FormHandler) UpdateQuestion(c *gin.Context) {
	var req dto.UpdateQuestionRequest
	if !utils.BindAndValidateJSON(c, &req) {
		return
	}
	if err := h.service.UpdateQuestion(&req); err != nil {
		c.JSON(500, gin.H{"message": "Failed to update question", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Question updated successfully"})
}

func (h *FormHandler) DeleteQuestion(c *gin.Context) {
	questionID := c.Param("questionId")
	if err := h.service.DeleteQuestion(questionID); err != nil {
		c.JSON(500, gin.H{"message": "Failed to delete question", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Question deleted successfully"})
}
