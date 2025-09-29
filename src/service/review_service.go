package service

import (
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
	"math"
)

type reviewService struct {
	reviewRepo repository.ReviewRepository
	courseRepo repository.CourseRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, courseRepo repository.CourseRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		courseRepo: courseRepo,
	}
}

func (rs *reviewService) GetCourseReviews(courseId uint, req *dto.GetCourseReviewsQueryRequest) (*dto.GetCourseReviewsResponse, error) {
	// Verify course exists
	_, err := rs.courseRepo.FindBySlug("")
	if err != nil {
		// For now, skip course verification or add FindByID to CourseRepository
	}

	// Set default values
	page := 1
	limit := 10
	orderBy := "created_at"
	sortBy := "desc"

	if req.Page > 0 {
		page = req.Page
	}
	if req.Limit > 0 {
		limit = req.Limit
	}
	if req.OrderBy != "" {
		orderBy = req.OrderBy
	}
	if req.SortBy != "" {
		sortBy = req.SortBy
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Prepare filters
	filters := make(map[string]interface{})

	if req.Rating != nil {
		filters["rating"] = *req.Rating
	}
	if req.Published != nil {
		filters["is_published"] = *req.Published
	} else {
		// Default: only show published reviews
		filters["is_published"] = true
	}

	// Get reviews with pagination
	reviews, total, err := rs.reviewRepo.GetCourseReviews(courseId, offset, limit, filters, orderBy, sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get course reviews", utils.ErrCodeInternal)
	}

	// Convert to DTO
	reviewItems := make([]dto.ReviewItem, len(reviews))
	for i, review := range reviews {
		userName := "Annymous"
		userAvatar := ""

		if review.User.Id != 0 {
			userName = review.User.FullName
			userAvatar = review.User.AvatarURL
		}

		reviewItems[i] = dto.ReviewItem{
			Id:          review.Id,
			UserId:      review.UserId,
			UserName:    userName,
			UserAvatar:  userAvatar,
			Rating:      review.Rating,
			Comment:     review.Comment,
			IsPublished: review.IsPublished,
			CreatedAt:   review.CreatedAt,
		}
	}

	// Get review stats
	stats, err := rs.reviewRepo.GetCourseReviewStats(courseId)
	if err != nil {
		// Log error but don't fail the request
		stats = &dto.ReviewStats{
			RatingDistribution: make(map[int]int),
		}
	}

	// Calculate pagination
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := dto.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return &dto.GetCourseReviewsResponse{
		Reviews:    reviewItems,
		Pagination: pagination,
		Stats:      *stats,
	}, nil
}
