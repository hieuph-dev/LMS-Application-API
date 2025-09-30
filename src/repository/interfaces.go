package repository

import (
	"lms/src/dto"
	"lms/src/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, bool)
	FindByUsername(username string) (*models.User, bool)
	FindById(id uint) (*models.User, error)
	UpdatePassword(userId uint, hashedPassword string) error
	UpdateProfile(userId uint, updates map[string]interface{}) error
	ChangePassword(userId uint, hashedPassword string) error
	UpdateAvatar(userId uint, avatarURL string) error
	GetUsersWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.User, int, error)
	DeleteUser(userId uint) error
}

type PasswordResetRepository interface {
	Create(reset *models.PasswordReset) error
	FindByToken(token string) (*models.PasswordReset, error)
	MarkAsUsed(id uint) error
	DeleteExpired() error
	DeleteByEmail(email string) error
}

type CategoryRepository interface {
	GetCategories(filters map[string]interface{}) ([]models.Category, int, error)
	FindById(id uint) (*models.Category, error)
	Create(category *models.Category) error
	FindBySlug(slug string) (*models.Category, bool)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
	HasChildren(id uint) (bool, error)
	FindBySlugExcept(slug string, excludeId uint) (*models.Category, bool)
}

type CourseRepository interface {
	GetCoursesWithPagination(offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error)
	SearchCourses(query string, offset, limit int, filters map[string]interface{}, sortBy, order string) ([]models.Course, int, error)
	GetSearchFilters(query string) (*dto.SearchFilters, error)
	GetFeaturedCourses(limit int, filters map[string]interface{}) ([]models.Course, int, error)
	FindBySlug(slug string) (*models.Course, error)
}

type ReviewRepository interface {
	GetCourseReviews(courseId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Review, int, error)
	GetCourseReviewStats(courseId uint) (*dto.ReviewStats, error)
}

type LessonRepository interface {
	GetCourseLessons(courseId uint) ([]models.Lesson, error)
	CheckUserEnrollment(userId, courseId uint) (bool, error)
	GetLessonProgress(userId uint, lessonIds []uint) (map[uint]bool, error)
}

type CouponRepository interface {
	FindByCode(code string) (*models.Coupon, error)
	IncrementUsedCount(couponId uint) error
	IsValidCoupon(coupon *models.Coupon) bool
}

type OrderRepository interface {
	Create(order *models.Order) error
	FindByOrderCode(orderCode string) (*models.Order, error)
	UpdatePaymentStatus(orderId uint, status string) error
}

type EnrollmentRepository interface {
	Create(enrollment *models.Enrollment) error
	CheckEnrollment(userId, courseId uint) (*models.Enrollment, bool)
	GetUserEnrollments(userId uint, offset, limit int, filters map[string]interface{}) ([]models.Enrollment, int, error)
	CompleteEnrollment(enrollmentId uint) error
}

type InstructorRepository interface {
	GetInstructorCourses(instructorId uint, offset, limit int, filters map[string]interface{}, orderBy, sortBy string) ([]models.Course, int, error)
	CreateCourse(course *models.Course) error
	FindCourseBySlug(slug string) (*models.Course, bool)
	FindCourseById(courseId uint) (*models.Course, error)
	FindCourseByIdAndInstructor(courseId, instructorId uint) (*models.Course, error)
	UpdateCourse(courseId uint, updates map[string]interface{}) error
	DeleteCourse(courseId uint) error
	CountEnrollmentsByCourse(courseId uint) (int64, error)
}
