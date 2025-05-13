package routes

import (
	"server/internal/handlers"

	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SubmissionRoutes(r *gin.Engine, handler *handlers.SubmissionHandler) {
	form := r.Group("/api/v1/forms")

	form.POST("/:id/submissions", handler.SendFormSubmission)

	admin := form.Group("", middleware.AuthRequired(), middleware.RoleOnly("user", "admin"))
	admin.GET("/:id/submissions", handler.GetFormSubmissions)
	admin.GET("/:id/submissions/:sessionid", handler.GetSubmissionsResult)
}
