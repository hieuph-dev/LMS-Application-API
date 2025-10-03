package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"math"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo      repository.OrderRepository
	courseRepo     repository.CourseRepository
	couponRepo     repository.CouponRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	courseRepo repository.CourseRepository,
	couponRepo repository.CouponRepository,
	enrollmentRepo repository.EnrollmentRepository,

) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		courseRepo:     courseRepo,
		enrollmentRepo: enrollmentRepo,
		couponRepo:     couponRepo,
	}
}

func (os *orderService) CreateOrder(userId uint, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := os.courseRepo.FindById(req.CourseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra course status
	if course.Status != "published" {
		return nil, utils.NewError("Course is not available for purchase", utils.ErrCodeBadRequest)
	}

	// 3. Kiểm tra user đã mua course chưa
	if existingEnrollment, exists := os.enrollmentRepo.CheckEnrollment(userId, req.CourseId); exists {
		if existingEnrollment.Status == "active" {
			return nil, utils.NewError("You already own this course", utils.ErrCodeConflict)
		}
	}

	// 4. Kiểm tra đã có order pending chưa
	existingOrder, err := os.orderRepo.FindPendingOrderByUserAndCourse(userId, req.CourseId)
	if err == nil && existingOrder != nil {
		return nil, utils.NewError("You already have a pending order for this course. Please complete or cancel it first", utils.ErrCodeConflict)
	}

	// 5. Tính giá gốc
	originalPrice := course.Price
	if course.DiscountPrice != nil && *course.DiscountPrice < originalPrice {
		originalPrice = *course.DiscountPrice
	}

	discountAmount := 0.0
	var couponId *uint
	var appliedCouponCode string

	// 6. Áp dụng coupon nếu có
	if req.CouponCode != "" {
		coupon, err := os.couponRepo.FindByCode(req.CouponCode)
		if err != nil {
			return nil, utils.NewError("Invalid coupon code", utils.ErrCodeBadRequest)
		}

		if !os.couponRepo.IsValidCoupon(coupon) {
			return nil, utils.NewError("Coupon is expired or not available", utils.ErrCodeBadRequest)
		}

		// Kiểm tra minimum order amount
		if originalPrice < coupon.MinOrderAmount {
			return nil, utils.NewError(
				fmt.Sprintf("Minimum order amount for this coupon is %2.f", coupon.MinOrderAmount),
				utils.ErrCodeBadRequest,
			)
		}

		// Tính discount
		if coupon.DiscountType == "percentage" {
			discountAmount = originalPrice * (coupon.DiscountValue / 100)
		} else if coupon.DiscountType == "fixed" {
			discountAmount = coupon.DiscountValue
		}

		// Apply max discount nếu có
		if coupon.MaxDiscountAmount != nil && discountAmount > *coupon.MaxDiscountAmount {
			discountAmount = *coupon.MaxDiscountAmount
		}

		couponId = &coupon.Id
		appliedCouponCode = coupon.Code
	}

	// 7. Tính final price
	finalPrice := originalPrice - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	// 8. Tạo order code
	orderCode := fmt.Sprintf("ORD-%s-%d", uuid.New().String()[:8], time.Now().Unix())

	// 9. Tạo order
	order := &models.Order{
		UserId:         userId,
		CourseId:       req.CourseId,
		OrderCode:      orderCode,
		OriginalPrice:  originalPrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponId:       couponId,
		PaymentStatus:  "pending",
	}

	if err := os.orderRepo.Create(order); err != nil {
		return nil, utils.WrapError(err, "Failed to create order", utils.ErrCodeInternal)
	}

	// 10. Nếu free course, tự động approve và tạo enrollment
	message := "Order created successfully. Please proceed to payment"
	if finalPrice == 0 {
		if err := os.completeOrder(order, "free"); err != nil {
			return nil, err
		}
		message = "Congratulations! You have successfully enrolled in this free course"
	}

	return &dto.CreateOrderResponse{
		OrderId:        order.Id,
		OrderCode:      order.OrderCode,
		CourseId:       course.Id,
		CourseTitle:    course.Title,
		OriginalPrice:  originalPrice,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		CouponCode:     appliedCouponCode,
		PaymentStatus:  order.PaymentStatus,
		CreatedAt:      order.CreatedAt,
		Message:        message,
	}, nil
}

func (os *orderService) GetOrderHistory(userId uint, req *dto.GetOrderHistoryQueryRequest) (*dto.GetOrderHistoryResponse, error) {
	// Set defaults
	page := 1
	limit := 10
	sortBy := "desc"

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}
	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})
	if req.PaymentStatus != "" {
		filters["payment_status"] = req.PaymentStatus
	}

	// Get orders
	orders, total, err := os.orderRepo.GetUsersOrders(userId, offset, limit, filters, "created_at", sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get order history", utils.ErrCodeInternal)
	}

	// Convert to DTO
	orderItems := make([]dto.OrderHistoryItem, len(orders))
	for i, order := range orders {
		// Load course info
		course, err := os.courseRepo.FindById(order.CourseId)
		if err != nil {
			orderItems[i] = dto.OrderHistoryItem{
				Id:              order.Id,
				OrderCode:       order.OrderCode,
				CourseId:        order.CourseId,
				CourseTitle:     "Course not found",
				CourseThumbnail: "",
				OriginalPrice:   order.OriginalPrice,
				DiscountAmount:  order.DiscountAmount,
				FinalPrice:      order.FinalPrice,
				PaymentStatus:   order.PaymentStatus,
				PaidAt:          order.PaidAt,
				CreatedAt:       order.CreatedAt,
			}
			continue
		}

		orderItems[i] = dto.OrderHistoryItem{
			Id:              order.Id,
			OrderCode:       order.OrderCode,
			CourseId:        order.CourseId,
			CourseTitle:     course.Title,
			CourseThumbnail: course.ThumbnailURL,
			OriginalPrice:   order.OriginalPrice,
			DiscountAmount:  order.DiscountAmount,
			FinalPrice:      order.FinalPrice,
			PaymentStatus:   order.PaymentStatus,
			PaidAt:          order.PaidAt,
			CreatedAt:       order.CreatedAt,
		}
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := dto.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &dto.GetOrderHistoryResponse{
		Orders:     orderItems,
		Pagination: pagination,
	}, nil
}

// Helper function to complete order and create enrollment
func (os *orderService) completeOrder(order *models.Order, paymentMethod string) error {
	now := time.Now()
	order.PaymentStatus = "paid"
	order.PaymentMethod = paymentMethod
	order.PaidAt = &now

	// Update order
	if err := os.orderRepo.Update(order); err != nil {
		return utils.WrapError(err, "Failed to update order", utils.ErrCodeInternal)
	}

	// Create enrollment
	enrollment := &models.Enrollment{
		UserId:             order.UserId,
		CourseId:           order.CourseId,
		EnrolledAt:         now,
		ProgressPercentage: 0,
		Status:             "active",
	}

	if err := os.enrollmentRepo.Create(enrollment); err != nil {
		return utils.WrapError(err, "Failed to create enrollment", utils.ErrCodeInternal)
	}

	// Update coupon used count nếu có
	if order.CouponId != nil {
		os.couponRepo.IncrementUsedCount(*order.CouponId)
	}

	return nil
}

func (os *orderService) GetOrderDetail(userId uint, orderId uint) (*dto.OrderDetailResponse, error) {
	// Tìm order
	order, err := os.orderRepo.FindById(orderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// Kiểm tra order có thuộc về user không
	if order.UserId != userId {
		return nil, utils.NewError("Access denied", utils.ErrCodeForbidden)
	}

	// Load course info
	course, err := os.courseRepo.FindById(order.CourseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// Get coupon code nếu có
	couponCode := ""
	if order.CouponId != nil {
		if coupon, err := os.couponRepo.FindById(*order.CouponId); err == nil {
			couponCode = coupon.Code
		}
	}

	return &dto.OrderDetailResponse{
		Id:              order.Id,
		OrderCode:       order.OrderCode,
		UserId:          order.UserId,
		CourseId:        order.CourseId,
		CourseTitle:     course.Title,
		CourseThumbnail: course.ThumbnailURL,
		InstructorName:  course.Instructor.FullName,
		OriginalPrice:   order.OriginalPrice,
		DiscountAmount:  order.DiscountAmount,
		FinalPrice:      order.FinalPrice,
		CouponCode:      couponCode,
		PaymentStatus:   order.PaymentStatus,
		PaidAt:          order.PaidAt,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}, nil
}

func (os *orderService) PayOrder(userId uint, orderId uint, req *dto.PayOrderRequest) (*dto.PayOrderResponse, error) {
	// 1. Tìm order
	order, err := os.orderRepo.FindById(orderId)
	if err != nil {
		return nil, utils.NewError("Order not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra order có thuộc về user không
	if order.UserId != userId {
		return nil, utils.NewError("Access denied", utils.ErrCodeForbidden)
	}

	// 3. Kiểm tra order status
	if order.PaymentStatus != "pending" {
		return nil, utils.NewError("Order has already been processed", utils.ErrCodeBadRequest)
	}

	// 4. Kiểm tra nếu là free course
	if order.FinalPrice == 0 {
		return nil, utils.NewError("This is a free order, no payment required", utils.ErrCodeBadRequest)
	}

	// 5. Simulate payment processing
	// TODO: Integrate with real payment gateway
	time.Sleep(1 * time.Second) // Simulate payment processing

	// 6. Complete order
	order.PaymentMethod = req.PaymentMethod
	if err := os.completeOrder(order, req.PaymentMethod); err != nil {
		return nil, err
	}

	return &dto.PayOrderResponse{
		OrderId:       order.Id,
		OrderCode:     order.OrderCode,
		PaymentStatus: "paid",
		PaymentMethod: req.PaymentMethod,
		PaidAt:        *order.PaidAt,
		Message:       "Payment successful! You have been enrolled in the course",
	}, nil
}
