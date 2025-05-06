package service

import (
	"github.com/yoanesber/Go-Department-CRUD/model"
	"github.com/yoanesber/Go-Department-CRUD/repository"
)

// Interface for department service
// This interface defines the methods that the department service should implement
type DepartmentService interface {
	GetAllDepartments() ([]model.Department, error)
	GetDepartmentByID(id string) (model.Department, error)
	CreateDepartment(department model.Department) (model.Department, error)
	UpdateDepartment(id string, department model.Department) (model.Department, error)
	DeleteDepartment(id string) (bool, error)
}

// This struct defines the DepartmentService that contains a repository field of type DepartmentRepository
type departmentService struct {
	repo repository.DepartmentRepository
}

// NewDepartmentService creates a new instance of DepartmentService with the given repository.
// It initializes the departmentService struct and returns it.
func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

// GetAllDepartments retrieves all departments from the database.
func (s *departmentService) GetAllDepartments() ([]model.Department, error) {
	departments, err := s.repo.GetAllDepartments()
	if err != nil {
		return nil, err
	}

	return departments, nil
}

// GetDepartmentByID retrieves a department by its ID from the database.
func (s *departmentService) GetDepartmentByID(id string) (model.Department, error) {
	department, err := s.repo.GetDepartmentByID(id)
	if err != nil {
		return model.Department{}, err
	}

	return department, nil
}

// CreateDepartment creates a new department in the database.
func (s *departmentService) CreateDepartment(department model.Department) (model.Department, error) {
	createdDepartment, err := s.repo.CreateDepartment(department)
	if err != nil {
		return model.Department{}, err
	}

	return createdDepartment, nil
}

// UpdateDepartment updates an existing department in the database.
func (s *departmentService) UpdateDepartment(id string, department model.Department) (model.Department, error) {
	updatedDepartment, err := s.repo.UpdateDepartment(id, department)
	if err != nil {
		return model.Department{}, err
	}

	return updatedDepartment, nil
}

// DeleteDepartment deletes a department by its ID from the database.
func (s *departmentService) DeleteDepartment(id string) (bool, error) {
	f, err := s.repo.DeleteDepartment(id)
	if err != nil {
		return false, err
	}

	return f, nil
}
