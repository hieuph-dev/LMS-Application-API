package routes

import (
	"lms/src/handler"

	"github.com/gin-gonic/gin"
)

type CourseRoutes struct {
	handler *handler.CourseHandler
}

func NewCourseRoutes(handler *handler.CourseHandler) *CourseRoutes {
	return &CourseRoutes{
		handler: handler,
	}
}

func (cr *CourseRoutes) Register(r *gin.RouterGroup) {
	courses := r.Group("/courses")
	{
		// Public routes
		courses.GET("/", cr.handler.GetCourses)
		courses.GET("/search", cr.handler.SearchCourses)
		courses.GET("/featured", cr.handler.GetFeaturedCourses)
		courses.GET("/:slug", cr.handler.GetCourseBySlug)
		courses.GET("/course_id/:course_id/reviews", cr.handler.GetCourseReviews)
	}
}
