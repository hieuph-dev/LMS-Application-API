package repository

import "gorm.io/gorm"

type DBLessonRepository struct {
	db *gorm.DB
}

func NewDBLessonRepository(db *gorm.DB) LessonRepository {
	return &DBLessonRepository{
		db: db,
	}
}
