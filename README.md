# RESTful API for Learning Management System (LMS)

A RESTful API built with Go (Golang) and Gin framework for an online learning platform, supporting course management, user roles, payments, and progress tracking.

## Key Features

- **Authentication & Authorization**: JWT-based login/register, password reset, role-based access (Admin, Instructor, Student).
- **User Management**: Profile updates, avatar upload, password management, user analytics.
- **Course Management**: CRUD operations, categorization, levels, search, ratings, and reviews.
- **Lessons**: CRUD, video lessons, ordering, previews.
- **Enrollment & Progress**: Course enrollment, progress tracking, completion certificates.
- **Payments & Orders**: Order creation, coupons, multiple payment methods, order history.
- **Coupons**: Discount types, validation rules, usage limits.
- **Analytics**: Revenue, student, course, and enrollment analytics for instructors and admins.

## Technology Stack

- **Go**: High-performance language for scalable APIs.
- **Gin**: Lightweight web framework for Go.
- **PostgreSQL**: Robust relational database with GORM ORM.
- **JWT**: Secure authentication with golang-jwt/jwt.
- **bcrypt**: Password hashing for security.

## API Documentation

Explore the full range of API endpoints and their usage in the interactive Postman documentation:

- **Postman documentation:** Link to Postman documentation:
https://documenter.getpostman.com/view/19784956/2sB3QJMAXQ

## Middleware

- **Auth**: Verifies JWT and sets user context.
- **Admin/Instructor**: Role-based access checks.
- **Rate Limiter**: 5 req/s, burst 10.
- **Logger**: Logs request/response details.

## Database Models

- **User**: Info, role, status, email verification.
- **Course**: Title, pricing, metadata, stats.
- **Lesson**: Title, video, order, publish status.
- **Enrollment**: User-course relation, progress, status.
- **Order**: Transaction, payment, coupon details.
- **Progress**: Lesson completion, watch duration.
- **Review**: Rating, comment, status.
- **Coupon**: Discount type, validation rules.

## Security

- Password hashing (bcrypt)
- JWT with expiration
- Role-based access
- Rate limiting
- SQL injection/XSS prevention

## Error Handling

Custom error codes (400, 401, 403, 404, 409, 500) with JSON response format.

## Performance

- Connection pooling (max 50)
- Optimized queries with GORM
- Indexed fields
- Pagination
- Rate limiting

## Get Started

1. **Clone the repository**:
    
    ```bash
    git clone https://github.com/hieuph-dev/LMS-Application-API.git
    cd lms
    
    ```
    
2. **Install dependencies**:
    
    ```bash
    go mod download
    
    ```
    
3. **Configure environment**: Create `.env` file:
    
    ```
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=lms_db
    JWT_SECRET=your-secret-key
    API_KEY=your-api-key
    
    ```
    
4. **Create database**:
    
    ```bash
    createdb lms_db
    
    ```
    
5. **Run migrations**: Auto-run on first start.
6. **Start the server**:
    
    ```bash
    go run cmd/api/main.go
    
    ```
    
7. Explore the API using Postman or similar tools.
