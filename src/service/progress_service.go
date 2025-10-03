package service

import (
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
	"time"
)

type progressService struct {
	progressRepo   repository.ProgressRepository
	enrollmentRepo repository.EnrollmentRepository
	courseRepo     repository.CourseRepository
	lessonRepo     repository.LessonRepository
}

func NewProgressService(
	progressRepo repository.ProgressRepository,
	enrollmentRepo repository.EnrollmentRepository,
	courseRepo repository.CourseRepository,
	lessonRepo repository.LessonRepository,
) ProgressService {
	return &progressService{
		progressRepo:   progressRepo,
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		lessonRepo:     lessonRepo,
	}
}

func (ps *progressService) GetCourseProgress(userId, courseId uint) (*dto.GetCourseProgressResponse, error) {
	// 1. Kiểm tra course có tồn tại không
	course, err := ps.courseRepo.FindById(courseId)
	if err != nil {
		return nil, utils.NewError("Course not found", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra user đã enroll chưa
	enrollment, isEnrolled := ps.enrollmentRepo.CheckEnrollment(userId, courseId)
	if !isEnrolled {
		return nil, utils.NewError("You are not enrolled in this course", utils.ErrCodeForbidden)
	}

	// 3. Lấy danh sách lessons của course
	lessons, err := ps.lessonRepo.GetCourseLessons(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get course lessons", utils.ErrCodeInternal)
	}

	// 4. Lấy progress của tất cả lessons
	progressMap := make(map[uint]*dto.LessonProgressItem)
	courseProgress, err := ps.progressRepo.GetCourseProgress(userId, courseId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get progress", utils.ErrCodeInternal)
	}

	// Map progress theo lesson_id
	progressDataMap := make(map[uint]struct {
		isCompleted   bool
		completedAt   *time.Time
		watchDuration int
		lastPosition  int
	})

	for _, p := range courseProgress {
		progressDataMap[p.LessonId] = struct {
			isCompleted   bool
			completedAt   *time.Time
			watchDuration int
			lastPosition  int
		}{
			isCompleted:   p.IsCompleted,
			completedAt:   p.CompletedAt,
			watchDuration: p.WatchDuration,
			lastPosition:  p.LastPosition,
		}
	}

	// 5. Tính toán progress cho từng lesson
	totalDuration := 0   // tổng thời lượng video của tất cả lessons.
	watchedDuration := 0 // tổng thời lượng mà user đã xem.
	completedCount := 0  // số bài học user đã hoàn thành.

	lessonItems := make([]dto.LessonProgressItem, 0, len(lessons))

	for _, lesson := range lessons {
		totalDuration += lesson.VideoDuration

		progressPercent := 0.0
		isCompleted := false
		var completedAt *time.Time = nil
		watchDuration := 0
		lastPosition := 0

		// Kiểm tra có progress không
		if p, exists := progressDataMap[lesson.Id]; exists {
			isCompleted = p.isCompleted
			completedAt = p.completedAt
			watchDuration = p.watchDuration
			lastPosition = p.lastPosition
			watchedDuration += watchDuration

			if isCompleted {
				completedCount++
				progressPercent = 100.0
			} else if lesson.VideoDuration > 0 {
				progressPercent = float64(watchDuration) / float64(lesson.VideoDuration) * 100
				if progressPercent > 100 {
					progressPercent = 100
				}
			}
		}

		lessonItems = append(lessonItems, dto.LessonProgressItem{
			LessonId:        lesson.Id,
			Title:           lesson.Title,
			Slug:            lesson.Slug,
			LessonOrder:     lesson.LessonOrder,
			VideoDuration:   lesson.VideoDuration,
			IsCompleted:     isCompleted,
			CompletedAt:     completedAt,
			WatchDuration:   watchDuration,
			LastPosition:    lastPosition,
			ProgressPercent: progressPercent,
		})

		progressMap[lesson.Id] = &lessonItems[len(lessonItems)-1]
	}

	// 6. Tính progress percentage tổng thể
	overallProgress := 0.0
	if len(lessons) > 0 {
		overallProgress = float64(completedCount) / float64(len(lessons)) * 100
	}

	// 7. Trả về response
	return &dto.GetCourseProgressResponse{
		CourseId:           course.Id,
		CourseTitle:        course.Title,
		IsEnrolled:         true,
		EnrolledAt:         &enrollment.EnrolledAt,
		ProgressPercentage: overallProgress,
		TotalLessons:       len(lessons),
		CompletedLessons:   completedCount,
		TotalDuration:      totalDuration,
		WatchedDuration:    watchedDuration,
		LastAccessedAt:     enrollment.LastAccessedAt,
		Status:             enrollment.Status,
		Lessons:            lessonItems,
	}, nil
}
