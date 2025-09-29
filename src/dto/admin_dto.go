package dto

import "time"

type GetUsersQueryRequest struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Role    string `form:"role" binding:"omitempty,oneof=admin student instructor guest"`
	Status  string `form:"status" binding:"omitempty,oneof=active inactive banned"`
	Search  string `form:"search" binding:"omitempty,search"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=created_at updated_at username email"`
	SortBy  string `form:"sort_by" binding:"omitempty,oneof=asc desc"`
}

type AdminUserItem struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GetUsersResponse struct {
	Users      []AdminUserItem `json:"users"`
	Pagination PaginationInfo  `json:"pagination"`
}

type PaginationInfo struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type AdminUserDetail struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Bio           string    `json:"bio"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	FullName      string `json:"full_name" binding:"omitempty,min=2,max=100"`
	Phone         string `json:"phone" binding:"omitempty,max=20"`
	Bio           string `json:"bio" binding:"omitempty,max=500"`
	AvatarURL     string `json:"avatar_url" binding:"omitempty,url"`
	Role          string `json:"role" binding:"omitempty,oneof=admin student instructor guest"`
	Status        string `json:"status" binding:"omitempty,oneof=active inactive banned"`
	EmailVerified bool   `json:"email_verified"`
}

type UpdateUserResponse struct {
	Id            uint      `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Bio           string    `json:"bio"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	Status        string    `json:"status"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
	UserId  uint   `json:"user_id"`
}

type ChangeUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive banned"`
	Reason string `json:"reason" binding:"omitempty,max=500"`
}

type ChangeUserStatusResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}
