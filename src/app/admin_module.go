package app

import (
	"lms/src/db"
	"lms/src/handler"
	"lms/src/repository"
	"lms/src/routes"
	"lms/src/service"
)

type AdminModule struct {
	routes routes.Route
}

func NewAdminModule() *AdminModule {
	// Tạo repository để tương tác với database
	userRepo := repository.NewDBUserRepository(db.DB)
	courseRepo := repository.NewDBCourseRepository(db.DB)

	// Tạo service chứa business logic
	adminService := service.NewAdminService(userRepo, courseRepo)

	// Tạo handler xử lý HTTP requests
	adminHandler := handler.NewAdminHandler(adminService)

	// Tạo routes định nghĩa các endpoint
	adminRoutes := routes.NewAdminRoutes(adminHandler)

	return &AdminModule{routes: adminRoutes}
}

func (am *AdminModule) Routes() routes.Route {
	return am.routes
}
