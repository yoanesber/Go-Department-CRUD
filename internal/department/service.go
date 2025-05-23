package department

import (
	"context"
	"errors"
	"fmt"

	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/metacontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"gorm.io/gorm"
)

// Interface for department service
// This interface defines the methods that the department service should implement
type DepartmentService interface {
	GetAllDepartments(ctx context.Context) ([]Department, error)
	GetDepartmentByID(ctx context.Context, id string) (Department, error)
	CreateDepartment(ctx context.Context, department Department) (Department, error)
	UpdateDepartment(ctx context.Context, id string, department Department) (Department, error)
	DeleteDepartment(ctx context.Context, id string) (bool, error)
}

// This struct defines the DepartmentService that contains a repository field of type DepartmentRepository
type departmentService struct {
	repo DepartmentRepository
}

// NewDepartmentService creates a new instance of DepartmentService with the given
// It initializes the departmentService struct and returns it.
func NewDepartmentService(repo DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

// GetAllDepartments retrieves all departments from the database.
func (s *departmentService) GetAllDepartments(ctx context.Context) ([]Department, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	// Retrieve all departments from the repository
	departments, err := s.repo.GetAllDepartments(db)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all departments: %v", err))
		return nil, err
	}

	return departments, nil
}

// GetDepartmentByID retrieves a department by its ID from the database.
func (s *departmentService) GetDepartmentByID(ctx context.Context, id string) (Department, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return Department{}, errors.New("database connection is nil")
	}

	// Retrieve the department by ID from the repository
	department, err := s.repo.GetDepartmentByID(db, id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get department by ID: %v", err))
		return Department{}, err
	}

	return department, nil
}

// CreateDepartment creates a new department in the database.
func (s *departmentService) CreateDepartment(ctx context.Context, d Department) (Department, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return Department{}, errors.New("database connection is nil")
	}

	// Validate the department struct using the validator
	if err := d.Validate(); err != nil {
		return Department{}, err
	}

	var createdDepartment Department
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the ID already exists
		existingDepartment, err := s.repo.GetDepartmentByID(db, d.ID)
		if (err == nil) || !(existingDepartment.Equals(&Department{})) {
			return errors.New("department with the same ID already exists")
		}

		// Check if the department name already exists
		existingDepartment, err = s.repo.GetDepartmentByName(db, d.DeptName)
		if err == nil || !(existingDepartment.Equals(&Department{})) {
			return errors.New("department with the same name already exists")
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(ctx)
		if !ok {
			return errors.New("missing user context")
		}

		// Create the department
		d.CreatedBy = &meta.UserID
		d.UpdatedBy = d.CreatedBy
		createdDepartment, err = s.repo.CreateDepartment(ctx, tx, d)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create department: %v", err))
		return Department{}, err
	}

	return createdDepartment, nil
}

// UpdateDepartment updates an existing department in the database.
func (s *departmentService) UpdateDepartment(ctx context.Context, id string, d Department) (Department, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return Department{}, errors.New("database connection is nil")
	}

	// Validate the department struct using the validator
	if err := d.Validate(); err != nil {
		return Department{}, err
	}

	var updatedDepartment Department
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the department exists
		existingDepartment, err := s.repo.GetDepartmentByID(db, id)
		if err != nil {
			return err
		}

		// Check if the existing department is empty
		if (existingDepartment.Equals(&Department{})) {
			return errors.New("department not found") // Department not found
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(ctx)
		if !ok {
			return errors.New("missing user context")
		}

		// Save the updated department
		existingDepartment.DeptName = d.DeptName
		existingDepartment.Active = d.Active
		existingDepartment.UpdatedBy = &meta.UserID
		updatedDepartment, err = s.repo.UpdateDepartment(ctx, tx, existingDepartment)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to update department: %v", err))
		return Department{}, err
	}

	return updatedDepartment, nil
}

// DeleteDepartment deletes a department by its ID from the database.
func (s *departmentService) DeleteDepartment(ctx context.Context, id string) (bool, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return false, errors.New("database connection is nil")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the department exists
		existingDepartment, err := s.repo.GetDepartmentByID(db, id)
		if err != nil {
			return err
		}

		// Check if the existing department is empty
		if (existingDepartment.Equals(&Department{})) {
			return errors.New("department not found") // Department not found
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(ctx)
		if !ok {
			return errors.New("missing user context")
		}

		// Delete the department
		err = s.repo.DeleteDepartment(ctx, tx, existingDepartment, &meta.UserID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete department: %v", err))
		return false, err
	}

	return true, nil
}
