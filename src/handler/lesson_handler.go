package handler

import (
	"lms/src/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(service service.LessonService) *LessonHandler {
	return &LessonHandler{
		service: service,
	}
}

// GET /api/v1/courses/:id/lessons - Lấy lessons của course (enrolled students only)
func (lh *LessonHandler) GetCourseLessons(ctx *gin.Context) {
	// Lấy course ID từ URL parameter
	courseIdParam := ctx.Param("id")
	if courseIdParam == "" {
		utils.ResponseError(ctx, utils.NewError("Course Id is required", utils.ErrCodeBadRequest))
		return
	}

	// Convert string to uint
	courseId, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		utils.ResponseError(ctx, utils.NewError("Invalid course Id format", utils.ErrCodeBadRequest))
	}

	// Lấy user ID từ context (đã được set bởi AuthMiddleware)
	userId, exists := ctx.Get("user_id")
	if !exists {
		utils.ResponseError(ctx, utils.NewError("User information not found", utils.ErrCodeUnauthorized))
		return
	}

	// Gọi service để lấy lessons
	response, err := lh.service.GetCourseLessons(userId.(uint), uint(courseId))
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}
