package repository

import (
	"fmt"
	"lms/src/dto"
	"lms/src/models"
	"strings"

	"gorm.io/gorm"
)

type DBReviewRepository struct {
	db *gorm.DB
}

func NewDBReviewRepository(db *gorm.DB) ReviewRepository {
	return &DBReviewRepository{
		db: db,
	}
}

func (rr *DBReviewRepository) GetCourseReviews(courseId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Review, int, error) {
	var reviews []models.Review
	var total int64

	query := rr.db.Model(&models.Review{}).
		Preload("User").
		Where("course_id = ? AND deleted_at IS NULL", courseId)

	// Apply filters
	for field, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Count total records
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
	if err := query.Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, int(total), nil
}

func (rr *DBReviewRepository) GetCourseReviewStats(courseId uint) (*dto.ReviewStats, error) {
	var stats dto.ReviewStats

	// Get total reviews and average rating
	result := rr.db.Model(&models.Review{}).
		Select("COUNT(*) as total_reviews, AVG(rating) as average_rating").
		Where("course_id = ? AND is_published = true AND deleted_at IS NULL", courseId).
		Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}

	// Get rating distribution
	var ratingDist []struct {
		Rating int `json:"rating"`
		Count  int `json:"count"`
	}

	err := rr.db.Model(&models.Review{}).
		Select("rating, COUNT(*) as count").
		Where("course_id = ? AND is_published = true AND deleted_at IS NULL", courseId).
		Group("rating").
		Scan(&ratingDist).Error

	if err != nil {
		return nil, err
	}

	// Convert to map
	stats.RatingDistribution = make(map[int]int)
	for _, dist := range ratingDist {
		stats.RatingDistribution[dist.Rating] = dist.Count
	}

	// Initialize missing ratings with 0
	for i := 1; i <= 5; i++ {
		if _, exists := stats.RatingDistribution[i]; !exists {
			stats.RatingDistribution[i] = 0
		}
	}

	return &stats, nil
}
