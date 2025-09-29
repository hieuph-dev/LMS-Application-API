package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type LessonRoutes struct {
	handler *handler.LessonHandler
}

func NewLessonRoutes(handler *handler.LessonHandler) *LessonRoutes {
	return &LessonRoutes{
		handler: handler,
	}
}

func (lr *LessonRoutes) Register(r *gin.RouterGroup) {
	courses := r.Group("/courses")
	{
		// Protected route - cần authentication
		courses.Use(middleware.AuthMiddleware())
		{
			courses.GET("/id/:id/lessons", lr.handler.GetCourseLessons)
		}
	}
}
