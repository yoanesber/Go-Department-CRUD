package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yoanesber/Go-Department-CRUD/internal/role"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/metacontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"gorm.io/gorm"
)

// Interface for user service
// This interface defines the methods that the user service should implement
type UserService interface {
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByUserName(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, id int64, user User) (User, error)
	UpdateLastLogin(ctx context.Context, id int64, lastLogin time.Time) (bool, error)
	// DeleteUser(id int64) (bool, error)
}

// This struct defines the UserService that contains a repository field of type UserRepository
// It implements the UserService interface and provides methods for user-related operations
type userService struct {
	repo UserRepository
}

// NewUserService creates a new instance of UserService with the given repository.
// It initializes the userService struct and returns it.
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// GetAllUsers retrieves all users from the database.
func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return nil, errors.New("database connection is nil")
	}

	// Retrieve all users from the repository
	users, err := s.repo.GetAllUsers(db)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get all users: %v", err))
		return nil, err
	}

	return users, nil
}

// GetUserByID retrieves a user by its ID from the database.
func (s *userService) GetUserByID(ctx context.Context, id int64) (User, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return User{}, errors.New("database connection is nil")
	}

	// Retrieve the user by ID from the repository
	user, err := s.repo.GetUserByID(db, id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user by ID: %v", err))
		return User{}, err
	}

	return user, nil
}

// GetUserByUserName retrieves a user by their username from the database.
func (s *userService) GetUserByUserName(ctx context.Context, username string) (User, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return User{}, errors.New("database connection is nil")
	}

	// Retrieve the user by username from the repository
	user, err := s.repo.GetUserByUserName(db, username)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user by username: %v", err))
		return User{}, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email from the database.
func (s *userService) GetUserByEmail(ctx context.Context, email string) (User, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return User{}, errors.New("database connection is nil")
	}

	// Retrieve the user by email from the repository
	user, err := s.repo.GetUserByEmail(db, email)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user by email: %v", err))
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user in the database.
func (s *userService) CreateUser(ctx context.Context, user User) (User, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return User{}, errors.New("database connection is nil")
	}

	// Validate the user struct using the validator
	if err := user.Validate(); err != nil {
		return User{}, err
	}

	// Validate the user's roles
	if len(user.Roles) == 0 {
		return User{}, errors.New("user must have at least one role")
	}
	for _, userRole := range user.Roles {
		if err := userRole.Validate(); err != nil {
			return User{}, err
		}
	}

	var createdUser User
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the user's roles are valid
		rRepo := role.NewRoleRepository()
		rServ := role.NewRoleService(rRepo)
		for i := range user.Roles {
			existingRole, err := rServ.GetRoleByName(ctx, user.Roles[i].Name)
			if err != nil {
				return err
			}
			if existingRole.Equals(&role.Role{}) {
				return errors.New("role with the given name not found")
			}

			// Assign/update the role ID in the user struct
			user.Roles[i].ID = existingRole.ID
		}

		// Check if the username already exists
		existingUser, err := s.repo.GetUserByUserName(db, user.UserName)
		if (err == nil) || !(existingUser.Equals(&User{})) {
			return errors.New("user with this username already exists")
		}

		// Check if the email already exists
		existingUser, err = s.repo.GetUserByEmail(db, user.Email)
		if (err == nil) || !(existingUser.Equals(&User{})) {
			return errors.New("user with this email already exists")
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(ctx)
		if !ok {
			return errors.New("missing user context")
		}

		// Create a new user in the database
		user.CreatedBy = &meta.UserID
		user.UpdatedBy = user.CreatedBy
		createdUser, err = s.repo.CreateUser(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create user: %v", err))
		return User{}, err
	}

	return createdUser, nil
}

// UpdateUser updates an existing user in the database.
func (s *userService) UpdateUser(ctx context.Context, id int64, user User) (User, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return User{}, errors.New("database connection is nil")
	}

	// Validate the user struct using the validator
	if err := user.Validate(); err != nil {
		return User{}, err
	}

	var updatedUser User
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		existingUser, err := s.repo.GetUserByID(db, id)
		if err != nil {
			return err
		}

		// Check if the existing user is empty
		if (existingUser.Equals(&User{})) {
			return errors.New("user not found") // User not found
		}

		// Extract user metadata from the context
		meta, ok := metacontext.ExtractRequestMeta(ctx)
		if !ok {
			return errors.New("missing user context")
		}

		// Update the user in the database
		existingUser.UserName = user.UserName
		existingUser.Password = user.Password
		existingUser.Email = user.Email
		existingUser.FirstName = user.FirstName
		existingUser.LastName = user.LastName
		existingUser.IsEnabled = user.IsEnabled
		existingUser.IsAccountNonExpired = user.IsAccountNonExpired
		existingUser.IsAccountNonLocked = user.IsAccountNonLocked
		existingUser.IsCredentialsNonExpired = user.IsCredentialsNonExpired
		existingUser.IsDeleted = user.IsDeleted
		existingUser.AccountExpirationDate = user.AccountExpirationDate
		existingUser.CredentialsExpirationDate = user.CredentialsExpirationDate
		existingUser.UserType = user.UserType
		existingUser.LastLogin = user.LastLogin
		existingUser.UpdatedBy = &meta.UserID
		existingUser.Roles = user.Roles
		updatedUser, err = s.repo.UpdateUser(ctx, tx, existingUser)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to update user: %v", err))
		return User{}, err
	}

	return updatedUser, nil
}

// UpdateLastLogin updates the last login time of a user in the database.
func (s *userService) UpdateLastLogin(ctx context.Context, id int64, lastLogin time.Time) (bool, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return false, errors.New("database connection is nil")
	}

	var isUpdated bool
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		existingUser, err := s.repo.GetUserByID(db, id)
		if err != nil {
			return err
		}

		// Check if the existing user is empty
		if (existingUser.Equals(&User{})) {
			return errors.New("user not found") // User not found
		}

		// Update the last login time
		*existingUser.LastLogin = lastLogin
		_, err = s.repo.UpdateUser(ctx, tx, existingUser)
		if err != nil {
			return err
		}

		// Set the isUpdated flag to true
		isUpdated = true

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to update last login: %v", err))
		return false, err
	}

	return isUpdated, nil
}
