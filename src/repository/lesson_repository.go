package repository

import (
	"lms/src/models"

	"gorm.io/gorm"
)

type DBLessonRepository struct {
	db *gorm.DB
}

func NewDBLessonRepository(db *gorm.DB) LessonRepository {
	return &DBLessonRepository{
		db: db,
	}
}

func (lr *DBLessonRepository) GetCourseLessons(courseId uint) ([]models.Lesson, error) {
	var lessons []models.Lesson

	err := lr.db.Where("course_id = ? AND is_published = ? AND deleted_at IS NULL", courseId, true).
		Order("lesson_order ASC").
		Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (lr *DBLessonRepository) CheckUserEnrollment(userId, courseId uint) (bool, error) {
	var count int64

	err := lr.db.Model(&models.Enrollment{}).
		Where("user_id = ? AND course_id = ? AND status = ? AND deleted_at IS NULL",
			userId, courseId, "active").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (lr *DBLessonRepository) GetLessonProgress(userId uint, lessonIds []uint) (map[uint]bool, error) {
	if len(lessonIds) == 0 {
		return make(map[uint]bool), nil
	}

	var progressList []models.Progress

	err := lr.db.Where("user_id = ? AND lesson_id IN ? AND deleted_at IS NULL",
		userId, lessonIds).
		Find(&progressList).Error

	if err != nil {
		return nil, err
	}

	progressMap := make(map[uint]bool)
	for _, progress := range progressList {
		progressMap[progress.LessonId] = progress.IsCompleted
	}

	return progressMap, nil
}
