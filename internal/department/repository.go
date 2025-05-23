package department

import (
	"context"
	"errors"

	"gorm.io/gorm" // Import GORM for ORM functionalities
)

// Interface for department repository
// This interface defines the methods that the department repository should implement
type DepartmentRepository interface {
	GetAllDepartments(tx *gorm.DB) ([]Department, error)
	GetDepartmentByID(tx *gorm.DB, id string) (Department, error)
	GetDepartmentByName(tx *gorm.DB, name string) (Department, error)
	CreateDepartment(ctx context.Context, tx *gorm.DB, d Department) (Department, error)
	UpdateDepartment(ctx context.Context, tx *gorm.DB, d Department) (Department, error)
	DeleteDepartment(ctx context.Context, tx *gorm.DB, d Department, deletedBy *int64) error
}

// This struct defines the DepartmentRepository that contains methods for interacting with the database
// It implements the DepartmentRepository interface and provides methods for department-related operations
type departmentRepository struct{}

// NewDepartmentRepository creates a new instance of DepartmentRepository.
// It initializes the departmentRepository struct and returns it.
func NewDepartmentRepository() DepartmentRepository {
	return &departmentRepository{}
}

// GetAllDepartments retrieves all departments from the database.
func (r *departmentRepository) GetAllDepartments(tx *gorm.DB) ([]Department, error) {
	var departments []Department
	err := tx.Order("id ASC").Find(&departments).Error
	if err != nil {
		return nil, err
	}

	return departments, nil
}

// It returns a slice of Department structs and an error if any occurs.
func (r *departmentRepository) GetDepartmentByID(tx *gorm.DB, id string) (Department, error) {
	var department Department
	err := tx.First(&department, "lower(id) = lower(?)", id).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, errors.New("department with the given ID not found")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, err
	}

	return department, nil
}

// GetDepartmentByName retrieves a department by its name from the database.
func (r *departmentRepository) GetDepartmentByName(tx *gorm.DB, name string) (Department, error) {
	var department Department
	err := tx.First(&department, "lower(dept_name) = lower(?)", name).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, errors.New("department with the given name not found")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return Department{}, err
	}

	return department, nil
}

// CreateDepartment inserts a new department into the database and returns the created department.
func (r *departmentRepository) CreateDepartment(ctx context.Context, tx *gorm.DB, d Department) (Department, error) {
	// Insert new department
	if err := tx.WithContext(ctx).Create(&d).Error; err != nil {
		return Department{}, err
	}

	return d, nil
}

// UpdateDepartment updates an existing department in the database and returns the updated department.
// It takes the department ID and the updated department struct as parameters.
func (r *departmentRepository) UpdateDepartment(ctx context.Context, tx *gorm.DB, d Department) (Department, error) {
	// Save the updated department
	if err := tx.WithContext(ctx).Save(&d).Error; err != nil {
		return Department{}, err
	}

	return d, nil
}

// DeleteDepartment deletes a department from the database by its ID.
// It takes the department ID as a parameter and returns an error if any occurs.
func (r *departmentRepository) DeleteDepartment(ctx context.Context, tx *gorm.DB, d Department, deletedBy *int64) error {
	// Set the deleted_by field to the user ID
	d.DeletedBy = deletedBy

	// Update the deleted_by field in the database
	// This is done to keep track of who deleted the department
	if err := tx.WithContext(ctx).Model(&d).Updates(Department{DeletedBy: deletedBy}).Error; err != nil {
		return err
	}

	// Delete the department from the database
	if err := tx.WithContext(ctx).Delete(&d).Error; err != nil {
		return err
	}

	return nil
}
