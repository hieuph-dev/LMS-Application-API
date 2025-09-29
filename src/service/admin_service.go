package service

import (
	"fmt"
	"lms/src/dto"
	"lms/src/repository"
	"lms/src/utils"
	"math"
	"strings"
	"time"
)

type adminService struct {
	userRepo repository.UserRepository
}

func NewAdminService(userRepo repository.UserRepository) AdminService {
	return &adminService{
		userRepo: userRepo,
	}
}

func (as *adminService) GetUsers(req *dto.GetUsersQueryRequest) (*dto.GetUsersResponse, error) {
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
	if req.Role != "" {
		filters["role"] = req.Role
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.Search != "" {
		filters["search"] = utils.NormalizeString(req.Search)
	}

	// Get users with pagination
	users, total, err := as.userRepo.GetUsersWithPagination(offset, limit, filters, orderBy, sortBy)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get users", utils.ErrCodeInternal)
	}

	// Convert to DTO
	userItems := make([]dto.AdminUserItem, len(users))
	for i, user := range users {
		userItems[i] = dto.AdminUserItem{
			Id:            user.Id,
			Username:      user.Username,
			Email:         user.Email,
			FullName:      user.FullName,
			Phone:         user.Phone,
			Role:          user.Role,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := dto.PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return &dto.GetUsersResponse{
		Users:      userItems,
		Pagination: pagination,
	}, nil
}

func (as *adminService) GetUserById(userId uint) (*dto.AdminUserDetail, error) {
	// Tìm user theo ID
	user, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// Convert sang DTO
	return &dto.AdminUserDetail{
		Id:            user.Id,
		Username:      user.Username,
		Email:         user.Email,
		FullName:      user.FullName,
		Phone:         user.Phone,
		Bio:           user.Bio,
		AvatarURL:     user.AvatarURL,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil
}

func (as *adminService) UpdateUser(userId uint, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	if existingUser.Status != "active" {
		return nil, utils.NewError("User account is not active", utils.ErrCodeForbidden)
	}

	// 2. Chuẩn bị dữ liệu cập nhật
	updates := make(map[string]interface{})

	if req.FullName != "" {
		updates["full_name"] = strings.TrimSpace(req.FullName)
	}
	if req.Phone != "" {
		updates["phone"] = strings.TrimSpace(req.Phone)
	}
	if req.Bio != "" {
		updates["bio"] = strings.TrimSpace(req.Bio)
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = strings.TrimSpace(req.AvatarURL)
	}
	if req.Role != "" {
		updates["role"] = strings.TrimSpace(req.Role)
	}
	if req.Status != "" {
		updates["status"] = strings.TrimSpace(req.Status)
	}
	// EmailVerified có thể là false nên kiểm tra khác
	updates["email_verified"] = req.EmailVerified

	updates["updated_at"] = time.Now()

	// 3. Cập nhật user
	if err := as.userRepo.UpdateProfile(userId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to update user", utils.ErrCodeInternal)
	}

	// 4. Lấy thông tin user đã cập nhật
	updatedUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to get updated user", utils.ErrCodeInternal)
	}

	return &dto.UpdateUserResponse{
		Id:            updatedUser.Id,
		Username:      updatedUser.Username,
		Email:         updatedUser.Email,
		FullName:      updatedUser.FullName,
		Phone:         updatedUser.Phone,
		Bio:           updatedUser.Bio,
		AvatarURL:     updatedUser.AvatarURL,
		Role:          updatedUser.Role,
		Status:        updatedUser.Status,
		EmailVerified: updatedUser.EmailVerified,
		CreatedAt:     updatedUser.CreatedAt,
		UpdatedAt:     updatedUser.UpdatedAt,
	}, nil
}

func (as *adminService) DeleteUser(userId uint) (*dto.DeleteUserResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 2. Không cho phép xóa admin khác
	if existingUser.Role == "admin" {
		return nil, utils.NewError("Cannot delete admin account", utils.ErrCodeForbidden)
	}

	// 3. Xóa user khỏi database (soft delete vì model có DeletedAt)
	if err := as.userRepo.DeleteUser(userId); err != nil {
		return nil, utils.WrapError(err, "Failed to delete user", utils.ErrCodeInternal)
	}

	return &dto.DeleteUserResponse{
		Message: "User deleted successfully",
		UserId:  userId,
	}, nil
}

func (as *adminService) ChangeUserStatus(userId uint, req *dto.ChangeUserStatusRequest) (*dto.ChangeUserStatusResponse, error) {
	// 1. Kiểm tra user có tồn tại không
	existingUser, err := as.userRepo.FindById(userId)
	if err != nil {
		return nil, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 2. Không cho phép thay đổi trạng thái admin khác
	if existingUser.Role == "admin" {
		return nil, utils.NewError("Cannot change admin account status", utils.ErrCodeForbidden)
	}

	// 3. Kiểm tra trạng thái hiện tại
	if existingUser.Status == req.Status {
		return nil, utils.NewError(fmt.Sprintf("User is already %s", req.Status), utils.ErrCodeBadRequest)
	}

	// 4. Cập nhật trạng thái
	updates := map[string]interface{}{
		"status":     req.Status,
		"updated_at": time.Now(),
	}

	if err := as.userRepo.UpdateProfile(userId, updates); err != nil {
		return nil, utils.WrapError(err, "Failed to updated user status", utils.ErrCodeInternal)
	}

	// 5. Tạo message tùy theo trạng thái
	var message string
	switch req.Status {
	case "active":
		message = "User account has been activated"
	case "inactive":
		message = "User account has been deactivated"
	case "banned":
		message = "User account has been banned"
	default:
		message = "User status has been updated"
	}

	if req.Reason != "" {
		message += fmt.Sprintf(". Reason: %s", req.Reason)
	}

	return &dto.ChangeUserStatusResponse{
		Id:       existingUser.Id,
		Username: existingUser.Username,
		Email:    existingUser.Email,
		Status:   req.Status,
		Message:  message,
	}, nil
}
