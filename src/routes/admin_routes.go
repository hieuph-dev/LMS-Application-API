package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct {
	handler *handler.AdminHandler
}

func NewAdminRoutes(handler *handler.AdminHandler) *AdminRoutes {
	return &AdminRoutes{
		handler: handler,
	}
}

func (ar *AdminRoutes) Register(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	{
		// All admin routes require authentication and admin role
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			// User management
			admin.GET("/users", ar.handler.GetUsers)
			admin.GET("/users/:id", ar.handler.GetUserById)
			admin.PUT("/users/:id", ar.handler.UpdateUser)
			admin.DELETE("/users/:id", ar.handler.DeleteUser)
			admin.PUT("/users/:id/status", ar.handler.ChangeUserStatus)
		}
	}
}
