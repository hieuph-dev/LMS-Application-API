package service

import (
	"lms/src/dto"
	"lms/src/models"
	"lms/src/repository"
	"lms/src/utils"
	"math"
)

type instructorService struct {
	instructorRepo repository.InstructorRepository
	categoryRepo   repository.CategoryRepository
}

func NewInstructorService(instructorRepo repository.InstructorRepository, categoryRepo repository.CategoryRepository) InstructorService {
	return &instructorService{
		instructorRepo: instructorRepo,
		categoryRepo:   categoryRepo,
	}
}

func (is *instructorService) GetInstructorCourses(instructorId uint, req *dto.GetInstructorCoursesQueryRequest) (*dto.GetInstructorCoursesResponse, error) {
	// Set default values
	page := 1
	if req.Page > 0 {
		page = req.Page
	}

	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := (page - 1) * limit

	// Build filters
	filters := make(map[string]interface{})

	if req.Status != "" {
		filters["status"] = req.Status
	}

	if req.Search != "" {
		filters["search"] = req.Search
	}

	// Get courses from repository
	courses, total, err := is.instructorRepo.GetInstructorCourses(
		instructorId,
		offset,
		limit,
		filters,
		req.OrderBy,
		req.SortBy,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to get instructor courses", utils.ErrCodeInternal)
	}

	// Convert to DTO
	courseItems := make([]dto.InstructorCourseItem, len(courses))
	for i, course := range courses {
		courseItems[i] = dto.InstructorCourseItem{
			Id:            course.Id,
			Title:         course.Title,
			Slug:          course.Slug,
			ThumbnailURL:  course.ThumbnailURL,
			Price:         course.Price,
			DiscountPrice: course.DiscountPrice,
			CategoryId:    course.CategoryId,
			CategoryName:  course.Category.Name,
			Level:         course.Level,
			Status:        course.Status,
			TotalLessons:  course.TotalLessons,
			DurationHours: course.DurationHours,
			EnrolledCount: course.EnrolledCount,
			RatingAvg:     course.RatingAvg,
			RatingCount:   course.RatingCount,
			IsFeatured:    course.IsFeatured,
			CreatedAt:     course.CreatedAt,
			UpdatedAt:     course.UpdatedAt,
		}
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.GetInstructorCoursesResponse{
		Courses: courseItems,
		Pagination: dto.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

func (is *instructorService) CreateCourse(instructorId uint, req *dto.CreateCourseRequest) (*dto.CreateCourseResponse, error) {
	// 1. Validate category exists và active
	category, err := is.categoryRepo.FindById(req.CategoryId)
	if err != nil {
		return nil, utils.NewError("category not found", utils.ErrCodeNotFound)
	}

	if !category.IsActive {
		return nil, utils.NewError("category is not active", utils.ErrCodeBadRequest)
	}

	// 2. Validate discount price
	if req.DiscountPrice != nil && *req.DiscountPrice >= req.Price {
		return nil, utils.NewError("discount price must be less than regular price", utils.ErrCodeBadRequest)
	}

	// 3. Generate unique slug
	baseSlug := utils.GenerateSlug(req.Title)
	uniqueSlug := utils.GenerateUniqueSlug(baseSlug, func(slug string) bool {
		_, exists := is.instructorRepo.FindCourseBySlug(slug)
		return exists
	})

	// 4. Create course model
	course := &models.Course{
		Title:         req.Title,
		Slug:          uniqueSlug,
		Description:   req.Description,
		ShortDesc:     req.ShortDesc,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		InstructorId:  instructorId,
		CategoryId:    req.CategoryId,
		Level:         req.Level,
		Language:      req.Language,
		Requirements:  req.Requirements,
		WhatYouLearn:  req.WhatYouLearn,
		DurationHours: req.DurationHours,
		Status:        "draft", // Mặc định là draft
		IsFeatured:    false,
		TotalLessons:  0,
		RatingAvg:     0,
		RatingCount:   0,
		EnrolledCount: 0,
	}

	// 5. Save to database
	if err := is.instructorRepo.CreateCourse(course); err != nil {
		return nil, utils.WrapError(err, "failed to create course", utils.ErrCodeInternal)
	}

	// 6. Return response
	return &dto.CreateCourseResponse{
		Id:              course.Id,
		Title:           course.Title,
		Slug:            course.Slug,
		Description:     course.Description,
		ShortDesc:       course.ShortDesc,
		ThumbnailURL:    course.ThumbnailURL,
		VideoPreviewURL: course.VideoPreviewURL,
		Price:           course.Price,
		DiscountPrice:   course.DiscountPrice,
		InstructorId:    course.InstructorId,
		CategoryId:      course.CategoryId,
		CategoryName:    category.Name,
		Level:           course.Level,
		DurationHours:   course.DurationHours,
		TotalLessons:    course.TotalLessons,
		Language:        course.Language,
		Requirements:    course.Requirements,
		WhatYouLearn:    course.WhatYouLearn,
		Status:          course.Status,
		IsFeatured:      course.IsFeatured,
		RatingAvg:       course.RatingAvg,
		RatingCount:     course.RatingCount,
		EnrolledCount:   course.EnrolledCount,
		CreatedAt:       course.CreatedAt,
		UpdatedAt:       course.UpdatedAt,
	}, nil
}

func (is *instructorService) UpdateCourse(instructorId, courseId uint, req *dto.UpdateCourseRequest) (*dto.UpdateCourseResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor này không
	course, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("course not found or you don't have permission to update this course", utils.ErrCodeNotFound)
	}

	// 2. Validate category nếu được cập nhật
	var category *models.Category
	if req.CategoryId != 0 && req.CategoryId != course.CategoryId {
		category, err = is.categoryRepo.FindById(req.CategoryId)
		if err != nil {
			return nil, utils.NewError("category not found", utils.ErrCodeNotFound)
		}
		if !category.IsActive {
			return nil, utils.NewError("category is not active", utils.ErrCodeBadRequest)
		}
	} else {
		// Giữ nguyên category hiện tại
		category, _ = is.categoryRepo.FindById(course.CategoryId)
	}

	// 3. Validate discount price nếu có
	finalPrice := course.Price
	if req.Price > 0 {
		finalPrice = req.Price
	}
	if req.DiscountPrice != nil && *req.DiscountPrice >= finalPrice {
		return nil, utils.NewError("discount price must be less than regular price", utils.ErrCodeBadRequest)
	}

	// 4. Build updates map
	updates := make(map[string]interface{})

	if req.Title != "" && req.Title != course.Title {
		updates["title"] = req.Title
		// Generate new slug nếu title thay đổi
		baseSlug := utils.GenerateSlug(req.Title)
		uniqueSlug := utils.GenerateUniqueSlug(baseSlug, func(slug string) bool {
			if slug == course.Slug {
				return false // Cho phép giữ nguyên slug hiện tại
			}
			_, exists := is.instructorRepo.FindCourseBySlug(slug)
			return exists
		})
		updates["slug"] = uniqueSlug
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.ShortDesc != "" {
		updates["short_desc"] = req.ShortDesc
	}

	if req.CategoryId != 0 {
		updates["category_id"] = req.CategoryId
	}

	if req.Level != "" {
		updates["level"] = req.Level
	}

	if req.Language != "" {
		updates["language"] = req.Language
	}

	if req.Price > 0 {
		updates["price"] = req.Price
	}

	if req.DiscountPrice != nil {
		updates["discount_price"] = req.DiscountPrice
	}

	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}

	if req.WhatYouLearn != "" {
		updates["what_you_learn"] = req.WhatYouLearn
	}

	if req.DurationHours >= 0 {
		updates["duration_hours"] = req.DurationHours
	}

	if req.Status != "" {
		// Không cho phép chuyển từ published về draft nếu đã có học viên
		if course.Status == "published" && req.Status == "draft" {
			enrollmentCount, _ := is.instructorRepo.CountEnrollmentsByCourse(courseId)
			if enrollmentCount > 0 {
				return nil, utils.NewError("cannot change status to draft when course has enrollments", utils.ErrCodeBadRequest)
			}
		}
		updates["status"] = req.Status
	}

	if req.IsFeatured != nil {
		updates["is_featured"] = *req.IsFeatured
	}

	// 5. Kiểm tra có gì cần update không
	if len(updates) == 0 {
		return nil, utils.NewError("no fields to update", utils.ErrCodeBadRequest)
	}

	// 6. Update course
	if err := is.instructorRepo.UpdateCourse(courseId, updates); err != nil {
		return nil, utils.WrapError(err, "failed to update course", utils.ErrCodeInternal)
	}

	// 7. Lấy lại course đã update
	updatedCourse, err := is.instructorRepo.FindCourseById(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch updated course", utils.ErrCodeInternal)
	}

	// 8. Return response
	return &dto.UpdateCourseResponse{
		Id:              updatedCourse.Id,
		Title:           updatedCourse.Title,
		Slug:            updatedCourse.Slug,
		Description:     updatedCourse.Description,
		ShortDesc:       updatedCourse.ShortDesc,
		ThumbnailURL:    updatedCourse.ThumbnailURL,
		VideoPreviewURL: updatedCourse.VideoPreviewURL,
		Price:           updatedCourse.Price,
		DiscountPrice:   updatedCourse.DiscountPrice,
		InstructorId:    updatedCourse.InstructorId,
		CategoryId:      updatedCourse.CategoryId,
		CategoryName:    category.Name,
		Level:           updatedCourse.Level,
		DurationHours:   updatedCourse.DurationHours,
		TotalLessons:    updatedCourse.TotalLessons,
		Language:        updatedCourse.Language,
		Requirements:    updatedCourse.Requirements,
		WhatYouLearn:    updatedCourse.WhatYouLearn,
		Status:          updatedCourse.Status,
		IsFeatured:      updatedCourse.IsFeatured,
		RatingAvg:       updatedCourse.RatingAvg,
		RatingCount:     updatedCourse.RatingCount,
		EnrolledCount:   updatedCourse.EnrolledCount,
		CreatedAt:       updatedCourse.CreatedAt,
		UpdatedAt:       updatedCourse.UpdatedAt,
	}, nil
}

func (is *instructorService) DeleteCourse(instructorId, courseId uint) (*dto.DeleteCourseResponse, error) {
	// 1. Kiểm tra course có tồn tại và thuộc về instructor này không
	course, err := is.instructorRepo.FindCourseByIdAndInstructor(courseId, instructorId)
	if err != nil {
		return nil, utils.NewError("course not found or you don't have permission to delete this course", utils.ErrCodeNotFound)
	}

	// 2. Kiểm tra xem course có học viên đang học không
	enrollmentCount, err := is.instructorRepo.CountEnrollmentsByCourse(courseId)
	if err != nil {
		return nil, utils.WrapError(err, "failed to check enrollments", utils.ErrCodeInternal)
	}

	if enrollmentCount > 0 {
		return nil, utils.NewError("cannot delete course with active enrollments", utils.ErrCodeBadRequest)
	}

	// 3. Chỉ cho phép xóa course ở trạng thái draft hoặc archived
	if course.Status == "published" {
		return nil, utils.NewError("cannot delete published course. Please archive it first", utils.ErrCodeBadRequest)
	}

	// 4. Xóa course (soft delete)
	if err := is.instructorRepo.DeleteCourse(courseId); err != nil {
		return nil, utils.WrapError(err, "failed to delete course", utils.ErrCodeInternal)
	}

	return &dto.DeleteCourseResponse{
		Message:  "Course deleted successfully",
		CourseId: courseId,
	}, nil
}
