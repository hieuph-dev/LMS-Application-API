package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type ProgressRoutes struct {
	handler *handler.ProgressHandler
}

func NewProgressRoutes(handler *handler.ProgressHandler) *ProgressRoutes {
	return &ProgressRoutes{
		handler: handler,
	}
}

func (pr *ProgressRoutes) Register(r *gin.RouterGroup) {
	enrollments := r.Group("/enrollments")
	{
		// Protected routes - cần authentication
		enrollments.Use(middleware.AuthMiddleware())
		{
			// Lấy learning progress của course
			enrollments.GET("/:course_id/progress", pr.handler.GetCourseProgress)
		}
	}
}
