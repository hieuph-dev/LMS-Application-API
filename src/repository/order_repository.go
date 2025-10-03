package repository

import (
	"fmt"
	"lms/src/models"
	"strings"

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

func (or *DBOrderRepository) Update(order *models.Order) error {
	return or.db.Save(order).Error
}

func (or *DBOrderRepository) FindById(orderId uint) (*models.Order, error) {
	var order models.Order
	if err := or.db.Where("id = ? AND deleted_at IS NULL", orderId).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *DBOrderRepository) FindByOrderCode(orderCode string) (*models.Order, error) {
	var order models.Order
	if err := or.db.Where("order_code = ? AND deleted_at IS NULL", orderCode).
		First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (or *DBOrderRepository) GetUsersOrders(userId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Order, int, error) {
	var orders []models.Order
	var total int64

	query := or.db.Model(&models.Order{}).
		Where("user_id = ? AND deleted_at IS NULL", userId)

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" && sortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", orderBy, strings.ToUpper(sortBy)))
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, int(total), nil
}

func (or *DBOrderRepository) FindPendingOrderByUserAndCourse(userId, courseId uint) (*models.Order, error) {
	var order models.Order
	err := or.db.Where("user_id = ? AND course_id = ? AND payment_status = ? AND deleted_at IS NULL",
		userId, courseId, "pending").
		First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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
