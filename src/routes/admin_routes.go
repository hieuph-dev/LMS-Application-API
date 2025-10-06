package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type AdminRoutes struct {
	handler       *handler.AdminHandler
	couponHandler *handler.CouponHandler // Thêm dòng này
}

func NewAdminRoutes(handler *handler.AdminHandler, couponHandler *handler.CouponHandler) *AdminRoutes {
	return &AdminRoutes{
		handler:       handler,
		couponHandler: couponHandler,
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

			// Course management
			admin.GET("/courses", ar.handler.GetCourses)
			admin.PUT("/courses/:course_id/status", ar.handler.ChangeCourseStatus)

			// Order management
			admin.GET("orders", ar.handler.GetAllOrders)
			admin.PUT("orders/:id/status", ar.handler.UpdateOrderStatus)

			// Coupon management - THÊM PHẦN NÀY
			admin.GET("/coupons", ar.couponHandler.GetAdminCoupons)
			admin.POST("/coupons", ar.couponHandler.CreateCoupon)
			admin.PUT("/coupons/:id", ar.couponHandler.UpdateCoupon)
			admin.DELETE("/coupons/:id", ar.couponHandler.DeleteCoupon)
		}
	}
}
