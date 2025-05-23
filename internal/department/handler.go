package department

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
	"gopkg.in/go-playground/validator.v9"
)

// This struct defines the DepartmentHandler which handles HTTP requests related to departments.
// It contains a service field of type DepartmentService which is used to interact with the department data layer.
type DepartmentHandler struct {
	Service DepartmentService
}

// NewDepartmentHandler creates a new instance of DepartmentHandler.
// It initializes the DepartmentHandler struct with the provided DepartmentService.
func NewDepartmentHandler(departmentService DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{Service: departmentService}
}

// GetAllDepartments retrieves all departments from the database and returns them as JSON.
// @Summary      Get all departments
// @Description  Get all departments from the database
// @Tags         departments
// @Accept       json
// @Produce      json
// @Success      200  {array}   HttpResponse for successful retrieval
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /departments [get]
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	departments, err := h.Service.GetAllDepartments(c.Request.Context())
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to retrieve departments", err.Error())
		return
	}

	util.JSONSuccess(c, http.StatusOK, "All Departments retrieved successfully", departments)
}

// GetDepartmentByID retrieves a department by its ID from the database and returns it as JSON.
// @Summary      Get department by ID
// @Description  Get a department by its ID from the database
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Department ID"
// @Success      200  {object}  HttpResponse for successful retrieval
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /departments/{id} [get]
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	// Parse the ID from the URL parameter
	id := c.Param("id")
	if id == "" {
		util.JSONError(c, http.StatusBadRequest, "Invalid ID", "ID cannot be empty")
		return
	}

	// Retrieve the department by ID from the service
	department, err := h.Service.GetDepartmentByID(c.Request.Context(), id)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to retrieve department", err.Error())
		return
	}

	if (department.Equals(&Department{})) {
		util.JSONError(c, http.StatusNotFound, "Department not found", "No department found with the given ID")
		return
	}

	util.JSONSuccess(c, http.StatusOK, "Department retrieved successfully", department)
}

// CreateDepartment creates a new department in the database and returns it as JSON.
// @Summary      Create a new department
// @Description  Create a new department in the database
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param        department  body      Department  true  "Department object"
// @Success      201  {object}  HttpResponse for successful creation
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	// Bind the JSON request body to the Department struct
	// and validate the input using ShouldBindJSON
	var department Department
	if err := c.ShouldBindJSON(&department); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Create the department using the service
	createdDepartment, err := h.Service.CreateDepartment(c.Request.Context(), department)
	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			util.JSONErrorMap(c, http.StatusBadRequest, "Failed to create department", util.FormatValidationErrors(err))
			return
		}

		util.JSONError(c, http.StatusInternalServerError, "Failed to create department", err.Error())
		return
	}

	util.JSONSuccess(c, http.StatusCreated, "Department created successfully", createdDepartment)
}

// UpdateDepartment updates an existing department in the database and returns it as JSON.
// @Summary      Update an existing department
// @Description  Update an existing department in the database
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param        id          path      string          true  "Department ID"
// @Param        department  body      Department  true  "Department object"
// @Success      200  {object}  HttpResponse for successful update
// @Failure      400  {object}  HttpResponse for bad request
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /departments/{id} [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	// Parse the ID from the URL parameter
	id := c.Param("id")
	if id == "" {
		util.JSONError(c, http.StatusBadRequest, "Invalid ID", "ID cannot be empty")
		return
	}

	// Bind the JSON request body to the Department struct
	var department Department
	if err := c.ShouldBindJSON(&department); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Update the department using the service
	department.ID = id // Set the ID of the department to be updated
	updatedDepartment, err := h.Service.UpdateDepartment(c.Request.Context(), id, department)
	if err != nil {
		// Check if the error is a validation error
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			util.JSONErrorMap(c, http.StatusBadRequest, "Failed to update department", util.FormatValidationErrors(err))
			return
		}

		util.JSONError(c, http.StatusInternalServerError, "Failed to update department", err.Error())
		return
	}

	// Check if the updated department is empty
	if (updatedDepartment.Equals(&Department{})) {
		util.JSONError(c, http.StatusNotFound, "Department not found", "No department found with the given ID")
		return
	}

	util.JSONSuccess(c, http.StatusOK, "Department updated successfully", updatedDepartment)
}

// DeleteDepartment deletes a department by its ID from the database.
// @Summary      Delete a department
// @Description  Delete a department by its ID from the database
// @Tags         departments
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Department ID"
// @Success      200  {object}  HttpResponse for successful deletion
// @Failure      404  {object}  HttpResponse for not found
// @Failure      500  {object}  HttpResponse for internal server error
// @Router       /departments/{id} [delete]
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	f, err := h.Service.DeleteDepartment(c.Request.Context(), id)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to delete department", err.Error())
		return
	}

	if !f {
		util.JSONError(c, http.StatusNotFound, "Department not found", "No department found with the given ID")
		return
	}

	util.JSONSuccess(c, http.StatusOK, "Department deleted successfully", nil)
}
