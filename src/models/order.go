package models

import (
	"time"

	"gorm.io/gorm"
)

// ---------------- Orders ----------------
type Order struct {
	Id             uint           `gorm:"primaryKey" json:"id"`
	UserId         uint           `json:"user_id"`
	CourseId       uint           `json:"course_id"`
	OrderCode      string         `gorm:"uniqueIndex;size:50;not null" json:"order_code"`
	OriginalPrice  float64        `gorm:"not null" json:"original_price"`
	DiscountAmount float64        `gorm:"default:0" json:"discount_amount"`
	FinalPrice     float64        `gorm:"not null" json:"final_price"`
	CouponId       *uint          `json:"coupon_id"`
	PaymentMethod  string         `gorm:"size:50" json:"payment_method"`
	PaymentStatus  string         `gorm:"size:20;default:pending" json:"payment_status"` // pending, paid, failed, refunded
	PaidAt         *time.Time     `json:"paid_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
