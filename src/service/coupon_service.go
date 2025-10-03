package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/repository"
)

type couponService struct {
	couponRepo repository.CouponRepository
	courseRepo repository.CourseRepository
}

func NewCouponService(
	couponRepo repository.CouponRepository,
	courseRepo repository.CourseRepository,
) CouponService {
	return &couponService{
		couponRepo: couponRepo,
		courseRepo: courseRepo,
	}
}

func (cs *couponService) ValidateCoupon(req *dto.ValidateCouponRequest) (*dto.ValidateCouponResponse, error) {
	// 1. Tìm coupon
	coupon, err := cs.couponRepo.FindByCode(req.CouponCode)
	if err != nil {
		return &dto.ValidateCouponResponse{
			Valid:   false,
			Message: "Invalid coupon code",
		}, nil
	}

	// 2. Kiểm tra coupon có valid không
	if !cs.couponRepo.IsValidCoupon(coupon) {
		return &dto.ValidateCouponResponse{
			Valid:      false,
			CouponCode: req.CouponCode,
			Message:    "Coupon is expired or not available",
		}, nil
	}

	// 3. Kiểm tra course có tồn tại không
	_, err = cs.courseRepo.FindById(req.CourseId)
	if err != nil {
		return &dto.ValidateCouponResponse{
			Valid:   false,
			Message: "Course not found",
		}, nil
	}

	// 4. Kiểm tra minimum order amount
	if req.OrderTotal < coupon.MinOrderAmount {
		return &dto.ValidateCouponResponse{
			Valid:          false,
			CouponCode:     req.CouponCode,
			MinOrderAmount: coupon.MinOrderAmount,
			Message:        fmt.Sprintf("Minimum order amount for this coupon is %.2f", coupon.MinOrderAmount),
		}, nil
	}

	// 5. Tính discount amount
	discountAmount := 0.0
	if coupon.DiscountType == "percentage" {
		discountAmount = req.OrderTotal * (coupon.DiscountValue / 100)
	} else if coupon.DiscountType == "fixed" {
		discountAmount = coupon.DiscountValue
	}

	// 6. Apply max discount nếu có
	if coupon.MaxDiscountAmount != nil && discountAmount > *coupon.MaxDiscountAmount {
		discountAmount = *coupon.MaxDiscountAmount
	}

	// 7. Tính final price
	finalPrice := req.OrderTotal - discountAmount
	if finalPrice < 0 {
		finalPrice = 0
	}

	return &dto.ValidateCouponResponse{
		Valid:             true,
		CouponCode:        coupon.Code,
		DiscountType:      coupon.DiscountType,
		DiscountValue:     coupon.DiscountValue,
		DiscountAmount:    discountAmount,
		FinalPrice:        finalPrice,
		MinOrderAmount:    coupon.MinOrderAmount,
		MaxDiscountAmount: coupon.MaxDiscountAmount,
		Message:           fmt.Sprintf("Coupon applied successfully! You save %.2f", discountAmount),
	}, nil
}
