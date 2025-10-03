package routes

import (
	"lms/src/handler"
	"lms/src/middleware"

	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	handler *handler.OrderHandler
}

func NewOrderRoutes(handler *handler.OrderHandler) *OrderRoutes {
	return &OrderRoutes{
		handler: handler,
	}
}

func (or *OrderRoutes) Register(r *gin.RouterGroup) {
	orders := r.Group("/orders")
	{
		// Protected routes - cần authentication
		orders.Use(middleware.AuthMiddleware())
		{
			// Create order
			orders.POST("/", or.handler.CreateOrder)

			// Get order history
			orders.GET("/", or.handler.GetOrderHistory)

			// Get order detail
			orders.GET("/:id", or.handler.GetOrderDetail)

			// Pay order
			orders.POST("/:id/pay", or.handler.PayOrder)
		}
	}

	// Coupon routes
	coupons := r.Group("/coupons")
	{
		coupons.Use(middleware.AuthMiddleware())
		{
			// Validate coupon
			coupons.POST("/validate", or.handler.ValidateCoupon)
		}

	}

}
