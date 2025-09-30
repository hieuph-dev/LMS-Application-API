package repository

import (
	"lms/src/models"
	"time"

	"gorm.io/gorm"
)

type DBCouponRepository struct {
	db *gorm.DB
}

func NewDBCouponRepository(db *gorm.DB) CouponRepository {
	return &DBCouponRepository{
		db: db,
	}
}

func (cr *DBCouponRepository) FindByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := cr.db.Where("code = ? AND is_active = ? AND deleted_at IS NULL", code, true).
		First(&coupon).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (cr *DBCouponRepository) IncrementUsedCount(couponId uint) error {
	return cr.db.Model(&models.Coupon{}).
		Where("id = ?", couponId).
		Update("used_count", gorm.Expr("used_count + 1")).Error
}

func (cr *DBCouponRepository) IsValidCoupon(coupon *models.Coupon) bool {
	now := time.Now()

	// Check if coupon is active
	if !coupon.IsActive {
		return false
	}

	// Check valid from date
	if coupon.ValidFrom != nil && now.Before(*coupon.ValidFrom) {
		return false
	}

	// Check valid to date
	if coupon.ValidTo != nil && now.After(*coupon.ValidTo) {
		return false
	}

	// Check usage limit
	if coupon.UsageLimit != nil && coupon.UsedCount >= *coupon.UsageLimit {
		return false
	}

	return true
}
