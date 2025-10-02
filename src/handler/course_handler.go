package handler

import (
	"lms/src/dto"
	"lms/src/service"
	"lms/src/utils"
	"lms/src/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service       service.CourseService
	reviewService service.ReviewService
}

func NewCourseHandler(service service.CourseService, reviewService service.ReviewService) *CourseHandler {
	return &CourseHandler{
		service:       service,
		reviewService: reviewService,
	}
}

// GET /api/v1/courses - Lấy danh sách courses (Public with filters and pagination)
func (ch *CourseHandler) GetCourses(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetCoursesQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy courses
	response, err := ch.service.GetCourses(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/courses/search - Search courses
func (ch *CourseHandler) SearchCourses(ctx *gin.Context) {
	// Parse query parameters
	var req dto.SearchCoursesQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để search courses
	response, err := ch.service.SearchCourses(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/courses/featured - Lấy courses nổi bật
func (ch *CourseHandler) GetFeaturedCourses(ctx *gin.Context) {
	// Parse query parameters
	var req dto.GetFeaturedCoursesQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy featured courses
	response, err := ch.service.GetFeaturedCourses(&req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

// GET /api/v1/courses/:slug - Lấy thông tin course detail
func (ch *CourseHandler) GetCourseBySlug(ctx *gin.Context) {
	// Lấy slug từ URL parameter
	slug := ctx.Param("slug")
	if slug == "" {
		utils.ResponseError(ctx, utils.NewError("Course slug is required", utils.ErrCodeBadRequest))
		return
	}

	// Gọi service để lấy thông tin course
	course, err := ch.service.GetCourseBySlug(slug)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, course)
}

// GET /api/v1/courses/course_id/:course_id/reviews - Lấy reviews của course
func (ch *CourseHandler) GetCourseReviews(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("course_id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
		return
	}

	// Parse query parameters
	var req dto.GetCourseReviewsQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ResponseValidator(ctx, validation.HandlerValidationErrors(err))
		return
	}

	// Gọi service để lấy reviews
	response, err := ch.reviewService.GetCourseReviews(uint(courseId), &req)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
