package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yoanesber/Go-Department-CRUD/model"
	"github.com/yoanesber/Go-Department-CRUD/service"
	"github.com/yoanesber/Go-Department-CRUD/util"
)

// This struct defines the DepartmentController which handles HTTP requests related to departments.
// It contains a service field of type DepartmentService which is used to interact with the department data layer.
type DepartmentController struct {
	Service service.DepartmentService
}

// GetAllDepartments retrieves all departments from the database and returns them as JSON.
// @Summary      Get all departments
// @Description  Get all departments from the database
// @Tags         departments
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.HttpResponse
// @Failure      500  {object}  model.HttpResponse
// @Router       /departments [get]
func (ctrl *DepartmentController) GetAllDepartments(c *gin.Context) {
	departments, err := ctrl.Service.GetAllDepartments()
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
// @Success      200  {object}  model.HttpResponse
// @Failure      404  {object}  model.HttpResponse
// @Failure      500  {object}  model.HttpResponse
// @Router       /departments/{id} [get]
func (ctrl *DepartmentController) GetDepartmentByID(c *gin.Context) {
	id := c.Param("id")
	department, err := ctrl.Service.GetDepartmentByID(id)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to retrieve department", err.Error())
		return
	}

	if (department == model.Department{}) {
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
// @Param        department  body      model.Department  true  "Department object"
// @Success      201  {object}  model.HttpResponse
// @Failure      400  {object}  model.HttpResponse
// @Failure      500  {object}  model.HttpResponse
// @Router       /departments [post]
func (ctrl *DepartmentController) CreateDepartment(c *gin.Context) {
	var department model.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	department.CreatedDate = time.Now()
	department.UpdatedDate = time.Now()
	department.UpdatedBy = department.CreatedBy

	createdDepartment, err := ctrl.Service.CreateDepartment(department)
	if err != nil {
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
// @Param        department  body      model.Department  true  "Department object"
// @Success      200  {object}  model.HttpResponse
// @Failure      400  {object}  model.HttpResponse
// @Failure      404  {object}  model.HttpResponse
// @Failure      500  {object}  model.HttpResponse
// @Router       /departments/{id} [put]
func (ctrl *DepartmentController) UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var department model.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		util.JSONError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	department.UpdatedDate = time.Now()

	updatedDepartment, err := ctrl.Service.UpdateDepartment(id, department)
	if err != nil {
		util.JSONError(c, http.StatusInternalServerError, "Failed to update department", err.Error())
		return
	}

	if (updatedDepartment == model.Department{}) {
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
// @Success      200  {object}  model.HttpResponse
// @Failure      404  {object}  model.HttpResponse
// @Failure      500  {object}  model.HttpResponse
// @Router       /departments/{id} [delete]
func (ctrl *DepartmentController) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	f, err := ctrl.Service.DeleteDepartment(id)
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
