package routes

import "lms/src/middleware"

type LessonRoutes struct {
	handler *handler.LessonHandler
}

func NewLessonRoutes(handler *handler.LessonHandler) *LessonRoutes {
	return &LessonRoutes{
		handler: handler,
	}
}

func (lr *LessonRoutes) Register(r *gin.RouteGroup) {
	courses := r.Group("/courses")
	{
		// Protected route - cần authentication
		courses.Use(middleware.AuthMiddleware())
		{
			courses.GET("/:id/lessons", lr.handler.GetCourseLessons)
		}
	}
}
