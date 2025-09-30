package service

import (
	"lms/src/dto"
	"mime/multipart"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(userId uint) (*dto.UserProfile, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.TokenResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
}

// Interface cho EmailService
type EmailService interface {
	SendPasswordResetEmail(email, resetToken, resetCode string) error
	SendWelcomeEmail(email, fullName string) error
}

type UserService interface {
	GetProfile(userId uint) (*dto.UserProfile, error)
	UpdateProfile(userId uint, req *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	ChangePassword(userId uint, req *dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
	UploadAvatar(userId uint, file *multipart.FileHeader) (*dto.UploadAvatarResponse, error)
}

type AdminService interface {
	GetUsers(req *dto.GetUsersQueryRequest) (*dto.GetUsersResponse, error)
	GetUserById(userId uint) (*dto.AdminUserDetail, error)
	UpdateUser(userId uint, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error)
	DeleteUser(userId uint) (*dto.DeleteUserResponse, error)
	ChangeUserStatus(userId uint, req *dto.ChangeUserStatusRequest) (*dto.ChangeUserStatusResponse, error)
}

type CategoryService interface {
	GetCategories(req *dto.GetCategoriesQueryRequest) (*dto.GetCategoriesResponse, error)
	GetCategoryById(categoryId uint) (*dto.CategoryDetail, error)
	CreateCategory(req *dto.CreateCategoryRequest) (*dto.CreateCategoryResponse, error)
	UpdateCategory(categoryId uint, req *dto.UpdateCategoryRequest) (*dto.UpdateCategoryResponse, error)
	DeleteCategory(categoryId uint) (*dto.DeleteCategoryResponse, error)
}

type CourseService interface {
	GetCourses(req *dto.GetCoursesQueryRequest) (*dto.GetCoursesResponse, error)
	SearchCourses(req *dto.SearchCoursesQueryRequest) (*dto.SearchCoursesResponse, error)
	GetFeaturedCourses(req *dto.GetFeaturedCoursesQueryRequest) (*dto.GetFeaturedCoursesResponse, error)
	GetCourseBySlug(slug string) (*dto.CourseDetail, error)
}

type ReviewService interface {
	GetCourseReviews(courseId uint, req *dto.GetCourseReviewsQueryRequest) (*dto.GetCourseReviewsResponse, error)
}

type LessonService interface {
	GetCourseLessons(userId, courseId uint) (*dto.GetCourseLessonsResponse, error)
}

type EnrollmentService interface {
	EnrollCourse(userId, courseId uint, req *dto.EnrollCourseRequest) (*dto.EnrollCourseResponse, error)
	CheckEnrollment(userId, courseId uint) (*dto.CheckEnrollmentResponse, error)
	GetMyEnrollments(userId uint, req *dto.GetMyEnrollmentsQueryRequest) (*dto.GetMyEnrollmentsResponse, error)
}
