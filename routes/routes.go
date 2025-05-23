package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/internal/auth"
	"github.com/yoanesber/Go-Department-CRUD/internal/dataredis"
	"github.com/yoanesber/Go-Department-CRUD/internal/department"
	"github.com/yoanesber/Go-Department-CRUD/internal/user"
	"github.com/yoanesber/Go-Department-CRUD/pkg/middleware/authorization"
	"github.com/yoanesber/Go-Department-CRUD/pkg/middleware/context"
	"github.com/yoanesber/Go-Department-CRUD/pkg/middleware/headers"
	"github.com/yoanesber/Go-Department-CRUD/pkg/middleware/logging"
	"github.com/yoanesber/Go-Department-CRUD/pkg/middleware/ratelimiter"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
	"golang.org/x/time/rate"
)

// SetupRouter initializes the router and sets up the routes for the application.
func SetupRouter() *gin.Engine {
	// Create a new Gin router instance
	r := gin.Default()

	// Set up middleware for the router
	// Middleware is used to handle cross-cutting concerns such as logging, security, and request ID generation
	r.Use(context.PostgresDBContext(), context.RedisContext(), headers.RequestSecurityHeader(), headers.RequestCorsHeader(),
		headers.RequestIDHeader(), logging.RequestLogger(), gzip.Gzip(gzip.DefaultCompression))

	// Set up the authentication routes
	// These routes handle user login and authentication
	authGroup := r.Group("/auth")
	{
		// Apply rate limiting middleware to the /auth group (e.g., login, register endpoints).
		// Configuration:
		// - Allows a burst of 1 request (no burst, basically one request at a time).
		// - After each request, only 1 new request is allowed every 30 seconds (refill rate).
		// - Each client IP has its own limiter instance which expires after 5 minutes of inactivity.
		authGroup.Use(ratelimiter.RateLimiter(rate.Every(30*time.Second), 1, 5*time.Minute))

		// Routes for authentication
		// These routes handle user login
		service := auth.NewAuthService()
		handler := auth.NewAuthHandler(service)

		// Define the routes for authentication
		// These routes handle user login
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/refresh-token", handler.RefreshToken)
	}

	// Set up the API version 1 routes
	v1 := r.Group("/api/v1", authorization.JwtValidation())
	{
		// Routes for department management
		// These routes handle CRUD operations for departments
		deptGroup := v1.Group("/departments")
		{
			// Apply rate limiting middleware to the /departments group.
			// Configuration:
			// - Allows up to 2 requests in quick succession (burst size = 2).
			// - After that, only 1 new request is allowed every 5 seconds (refill rate).
			// - Each client IP has its own limiter instance that expires after 10 minutes of inactivity.
			deptGroup.Use(ratelimiter.RateLimiter(rate.Every(5*time.Second), 2, 10*time.Minute))

			// Initialize the department repository and service
			// This is where the actual implementation of the repository and service would be used
			repo := department.NewDepartmentRepository()
			service := department.NewDepartmentService(repo)

			// Initialize the department handler with the service
			// This handler handles the HTTP requests and responses for department-related operations
			handler := department.NewDepartmentHandler(service)

			// Define the routes for department management
			// These routes handle CRUD operations for departments
			deptGroup.GET("", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), handler.GetAllDepartments)
			deptGroup.GET("/:id", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), handler.GetDepartmentByID)
			deptGroup.POST("", authorization.RoleBasedAccessControl("ROLE_ADMIN"), handler.CreateDepartment)
			deptGroup.PUT("/:id", authorization.RoleBasedAccessControl("ROLE_ADMIN"), handler.UpdateDepartment)
			deptGroup.DELETE("/:id", authorization.RoleBasedAccessControl("ROLE_ADMIN"), handler.DeleteDepartment)
		}

		// Routes for user management
		// These routes handle CRUD operations for users
		userGroup := v1.Group("/users")
		{
			// Rate limiter middleware for the /users group, accessible only by admin users.
			// - Allows a burst of up to 10 requests at once.
			// - Allows 1 request per second continuously after the burst.
			// - Limits each admin IP to prevent spamming the user management endpoints.
			// - Limiter TTL is 15 minutes to clean up inactive IP limiters.
			userGroup.Use(ratelimiter.RateLimiter(rate.Every(1*time.Second), 10, 15*time.Minute))

			// Initialize the user repository and service
			// This is where the actual implementation of the repository and service would be used
			repo := user.NewUserRepository()
			service := user.NewUserService(repo)

			// Initialize the user handler with the service
			// This handler handles the HTTP requests and responses for user-related operations
			handler := user.NewUserHandler(service)

			// Define the routes for user management
			// These routes handle CRUD operations for users
			userGroup.GET("", authorization.RoleBasedAccessControl("ROLE_ADMIN"), handler.GetAllUsers)
			userGroup.GET("/:id", authorization.RoleBasedAccessControl("ROLE_ADMIN"), handler.GetUserByID)
			userGroup.POST("", authorization.RoleBasedAccessControl("ROLE_ADMIN"), handler.CreateUser)
		}

		dataRedisGroup := v1.Group("/dataredis")
		{
			// Rate limiter middleware for the /dataredis group.
			// - Allows a burst of up to 5 requests at once.
			// - Allows 1 request every 3 seconds continuously after the burst.
			// - Helps prevent abuse of Redis storage/read operations from a single IP.
			// - Limiter TTL is 10 minutes to clean up inactive IP limiters.
			dataRedisGroup.Use(ratelimiter.RateLimiter(rate.Every(3*time.Second), 5, 10*time.Minute))

			// Initialize the data redis service
			// This is where the actual implementation of the service would be used
			service := dataredis.NewDataRedisService()

			// Initialize the data redis handler with the service
			// This handler handles the HTTP requests and responses for data redis-related operations
			handler := dataredis.NewDataRedisHandler(service)

			// Define the routes for data redis management
			dataRedisGroup.GET("/string/:key", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), handler.GetStringValue)
			dataRedisGroup.GET("/json/:key", authorization.RoleBasedAccessControl("ROLE_ADMIN", "ROLE_USER"), handler.GetJSONValue)
		}
	}

	// NoRoute handler for undefined routes
	// This handler will be called when no other route matches the request
	r.NoRoute(func(c *gin.Context) {
		util.JSONError(c, http.StatusNotFound, "Not Found", "The requested resource was not found")
	})

	// NoMethod handler for unsupported HTTP methods
	// This handler will be called when a request method is not allowed for the requested resource
	r.NoMethod(func(c *gin.Context) {
		util.JSONError(c, http.StatusMethodNotAllowed, "Method Not Allowed", "The requested method is not allowed for this resource")
	})

	return r
}
