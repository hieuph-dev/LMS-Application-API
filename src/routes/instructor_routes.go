package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type InstructorRoutes struct {
	handler *handler.InstructorHandler
}

func NewInstructorRoutes(handler *handler.InstructorHandler) *InstructorRoutes {
	return &InstructorRoutes{
		handler: handler,
	}
}

func (ir *InstructorRoutes) Register(r *gin.RouterGroup) {
	instructor := r.Group("/instructor")
	{
		// Protected routes - cần authentication và role instructor
		instructor.Use(middleware.AuthMiddleware())
		instructor.Use(middleware.InstructorMiddleware())
		{
			instructor.GET("/courses", ir.handler.GetInstructorCourses)
			instructor.POST("/courses", ir.handler.CreateCourse)
			instructor.PUT("/courses/:id", ir.handler.UpdateCourse)
			instructor.DELETE("/courses/:id", ir.handler.DeleteCourse)
			instructor.GET("/courses/:id/students", ir.handler.GetCourseStudents)
		}
	}
}
