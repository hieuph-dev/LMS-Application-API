package repository

import (
	"lms/src/models"

	"gorm.io/gorm"
)

type DBOrderRepository struct {
	db *gorm.DB
}

func NewDBOrderRepository(db *gorm.DB) OrderRepository {
	return &DBOrderRepository{
		db: db,
	}
}

func (or *DBOrderRepository) Create(order *models.Order) error {
	return or.db.Create(order).Error
}

func (or *DBOrderRepository) FindByOrderCode(orderCode string) (*models.Order, error) {
	var order models.Order
	err := or.db.Where("order_code = ?", orderCode).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *DBOrderRepository) UpdatePaymentStatus(orderId uint, status string) error {
	updates := map[string]interface{}{
		"payment_status": status,
	}

	if status == "paid" {
		updates["paid_at"] = gorm.Expr("NOW()")
	}

	return or.db.Model(&models.Order{}).
		Where("id = ?", orderId).
		Updates(updates).Error
}
