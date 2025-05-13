package routes

import (
	"server/internal/handlers"

	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func FormRoutes(r *gin.Engine, handler *handlers.FormHandler) {
	form := r.Group("/api/v1/forms", middleware.AuthRequired(), middleware.RoleOnly("user"))

	form.POST("", handler.CreateNewForm)
	form.GET("", handler.GetAllForms)
	form.GET("/:id", handler.GetFormDetail)

	form.GET("/:id/settings", handler.GetFormSettings)
	form.PUT("/:id/settings", handler.UpdateFormSettings)

	form.POST("/:id/sections", handler.AddFormSection)
	form.GET("/:id/sections", handler.GetFormSections)

	form.GET("/:id/questions", handler.GetFormQuestion)
	form.POST("/:id/questions", handler.AddFormQuestion)

	form.PUT("/sections/:sectionId", handler.UpdateFormSections)
	form.DELETE("/sections/:sectionId", handler.DeleteFormSections)

	form.PUT("/questions", handler.UpdateQuestion)
	form.DELETE("/questions", handler.DeleteQuestion)
}
