package repository

import (
	"fmt"
	"lms/src/models"
	"strings"

	"gorm.io/gorm"
)

type DBInstructorRepository struct {
	db *gorm.DB
}

func NewDBInstructorRepository(db *gorm.DB) InstructorRepository {
	return &DBInstructorRepository{
		db: db,
	}
}

func (ir *DBInstructorRepository) GetInstructorCourses(instructorId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error) {
	var courses []models.Course
	var total int64

	query := ir.db.Model(&models.Course{}).
		Preload("Category").
		Where("instructor_id = ?", instructorId)

	// Apply filters
	for field, value := range filters {
		if field == "search" {
			searchTerm := fmt.Sprintf("%%%s%%", value)
			query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		}
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
	if err := query.Offset(offset).Limit(limit).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, int(total), nil
}

func (ir *DBInstructorRepository) CreateCourse(course *models.Course) error {
	return ir.db.Create(course).Error
}

func (ir *DBInstructorRepository) FindCourseBySlug(slug string) (*models.Course, bool) {
	var course models.Course
	if err := ir.db.Where("slug = ?", slug).First(&course).Error; err != nil {
		return nil, false
	}
	return &course, true
}

func (ir *DBInstructorRepository) FindCourseById(courseId uint) (*models.Course, error) {
	var course models.Course
	if err := ir.db.Preload("Category").Where("id = ?", courseId).First(&course).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (ir *DBInstructorRepository) FindCourseByIdAndInstructor(courseId, instructorId uint) (*models.Course, error) {
	var course models.Course
	if err := ir.db.Preload("Category").
		Where("id = ? AND instructor_id = ?", courseId, instructorId).
		First(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (ir *DBInstructorRepository) UpdateCourse(courseId uint, updates map[string]interface{}) error {
	return ir.db.Model(&models.Course{}).Where("id = ?", courseId).Updates(updates).Error
}

func (ir *DBInstructorRepository) DeleteCourse(courseId uint) error {
	return ir.db.Delete(&models.Course{}, courseId).Error
}

func (ir *DBInstructorRepository) CountEnrollmentsByCourse(courseId uint) (int64, error) {
	var count int64
	if err := ir.db.Model(&models.Enrollment{}).
		Where("course_id = ?", courseId).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
