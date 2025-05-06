package repository

import (
	"errors"

	"github.com/yoanesber/Go-Department-CRUD/config"
	"github.com/yoanesber/Go-Department-CRUD/model"
)

// Interface for department repository
// This interface defines the methods that the department repository should implement
type DepartmentRepository interface {
	GetAllDepartments() ([]model.Department, error)
	GetDepartmentByID(id string) (model.Department, error)
	CreateDepartment(d model.Department) (model.Department, error)
	UpdateDepartment(id string, d model.Department) (model.Department, error)
	DeleteDepartment(id string) (bool, error)
}

// This struct defines the DepartmentRepository that contains methods for interacting with the database
type departmentRepository struct{}

// NewDepartmentRepository creates a new instance of DepartmentRepository.
// It initializes the departmentRepository struct and returns it.
func NewDepartmentRepository() DepartmentRepository {
	return &departmentRepository{}
}

// GetAllDepartments retrieves all departments from the database.
func (r *departmentRepository) GetAllDepartments() ([]model.Department, error) {
	var departments []model.Department
	_, err := config.GetDB().Select(&departments, "SELECT * FROM employees.department ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	return departments, nil
}

// It returns a slice of Department structs and an error if any occurs.
func (r *departmentRepository) GetDepartmentByID(id string) (model.Department, error) {
	var department model.Department
	err := config.GetDB().SelectOne(&department, "SELECT * FROM employees.department WHERE id = $1", id)
	if err != nil {
		return model.Department{}, err
	}

	return department, nil
}

// CreateDepartment inserts a new department into the database and returns the created department.
func (r *departmentRepository) CreateDepartment(d model.Department) (model.Department, error) {
	// Check if the department already exists
	existingDepartment, err := r.GetDepartmentByID(d.ID)
	if err != nil && err.Error() != "sql: no rows in result set" {
		// If the error is not "no rows", return the error
		return model.Department{}, err
	}

	// Check if the existing department is not empty
	// If it is not empty, it means the department already exists
	if (existingDepartment != model.Department{}) {
		return model.Department{}, errors.New("department already exists") // Department already exists
	}

	// Insert the new department into the database
	var returnId string
	err = config.GetDB().QueryRow(
		"INSERT INTO employees.department (id, dept_name, active, created_by, created_date, updated_by, updated_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		d.ID, d.DeptName, d.Active, d.CreatedBy, d.CreatedDate, d.UpdatedBy, d.UpdatedDate,
	).Scan(&returnId)
	if err != nil {
		return model.Department{}, err
	}

	// Check if the ID was returned
	// If the ID is empty, it means the insertion failed
	if returnId == "" {
		return model.Department{}, errors.New("failed to create department") // Failed to create department
	}

	return d, nil
}

// UpdateDepartment updates an existing department in the database and returns the updated department.
// It takes the department ID and the updated department struct as parameters.
func (r *departmentRepository) UpdateDepartment(id string, d model.Department) (model.Department, error) {
	// Check if the department exists
	existingDepartment, err := r.GetDepartmentByID(id)
	if err != nil {
		return model.Department{}, err
	}

	// Check if the existing department is empty
	if (existingDepartment == model.Department{}) {
		return model.Department{}, errors.New("department not found") // Department not found
	}

	// Update the existing department with the new values
	existingDepartment.DeptName = d.DeptName
	existingDepartment.Active = d.Active
	existingDepartment.UpdatedBy = d.UpdatedBy
	existingDepartment.UpdatedDate = d.UpdatedDate

	// Update the department object with the new values
	operation, err := config.GetDB().Exec(
		"UPDATE employees.department SET dept_name = $1, active = $2, updated_by = $3, updated_date = $4 WHERE id = $5",
		existingDepartment.DeptName, existingDepartment.Active, existingDepartment.UpdatedBy, existingDepartment.UpdatedDate, id,
	)
	if err != nil {
		return model.Department{}, err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return model.Department{}, errors.New("failed to update department") // Failed to update department
	}

	return existingDepartment, nil
}

// DeleteDepartment deletes a department from the database by its ID.
// It takes the department ID as a parameter and returns an error if any occurs.
func (r *departmentRepository) DeleteDepartment(id string) (bool, error) {
	// Check if the department exists
	existingDepartment, err := r.GetDepartmentByID(id)
	if err != nil {
		return false, err
	}

	// Check if the existing department is empty
	if (existingDepartment == model.Department{}) {
		return false, errors.New("department not found") // Department not found
	}

	// Delete the department from the database
	operation, err := config.GetDB().Exec("DELETE FROM employees.department WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return false, errors.New("failed to delete department") // Failed to delete department
	}

	return true, nil
}
