package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
)

// Interface for role service
// This interface defines the methods that the role service should implement
type RoleService interface {
	GetRoleByID(ctx context.Context, id uint) (Role, error)
	GetRoleByName(ctx context.Context, name string) (Role, error)
}

// This struct defines the RoleService that contains a repository field of type RoleRepository
// It implements the RoleService interface and provides methods for role-related operations
type roleService struct {
	repo RoleRepository
}

// NewRoleService creates a new instance of RoleService with the given repository.
// It initializes the roleService struct and returns it.
func NewRoleService(repo RoleRepository) RoleService {
	return &roleService{repo: repo}
}

// GetRoleByID retrieves a role by its ID from the database.
func (s *roleService) GetRoleByID(ctx context.Context, id uint) (Role, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return Role{}, errors.New("database connection is nil")
	}

	// Retrieve the role by ID from the repository
	role, err := s.repo.GetRoleByID(db, id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get role by ID: %v", err))
		return Role{}, err
	}

	return role, nil
}

// GetRoleByName retrieves a role by its name from the database.
func (s *roleService) GetRoleByName(ctx context.Context, name string) (Role, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return Role{}, errors.New("database connection is nil")
	}

	// Retrieve the role by name from the repository
	role, err := s.repo.GetRoleByName(db, name)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get role by name: %v", err))
		return Role{}, err
	}

	return role, nil
}
