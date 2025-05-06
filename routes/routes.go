package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/controller"
	"github.com/yoanesber/Go-Department-CRUD/repository"
	"github.com/yoanesber/Go-Department-CRUD/service"
	"github.com/yoanesber/Go-Department-CRUD/util"
)

// SetupRouter initializes the router and sets up the routes for the application.
func SetupRouter() *gin.Engine {
	// Create a new Gin router instance
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		// Routes for department management
		// These routes handle CRUD operations for departments
		dept := v1.Group("/departments")
		{
			// Initialize the department repository and service
			// This is where the actual implementation of the repository and service would be used
			repo := repository.NewDepartmentRepository()
			service := service.NewDepartmentService(repo)

			// Initialize the department controller with the service
			// This controller handles the HTTP requests and responses for department-related operations
			controller := &controller.DepartmentController{Service: service}

			// Define the routes for department management
			// These routes handle CRUD operations for departments
			dept.GET("", controller.GetAllDepartments)
			dept.GET("/:id", controller.GetDepartmentByID)
			dept.POST("", controller.CreateDepartment)
			dept.PUT("/:id", controller.UpdateDepartment)
			dept.DELETE("/:id", controller.DeleteDepartment)
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
