package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Interface for user repository
// This interface defines the methods that the user repository should implement
type UserRepository interface {
	GetAllUsers(tx *gorm.DB) ([]User, error)
	GetUserByID(tx *gorm.DB, id int64) (User, error)
	GetUserByUserName(tx *gorm.DB, username string) (User, error)
	GetUserByEmail(tx *gorm.DB, email string) (User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user User) (User, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user User) (User, error)
	// DeleteUser(id int64) (bool, error)
}

// This struct defines the UserRepository that contains methods for interacting with the database
// It implements the UserRepository interface and provides methods for user-related operations
type userRepository struct{}

// NewUserRepository creates a new instance of UserRepository.
// It initializes the userRepository struct and returns it.
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// GetAllUsers retrieves all users from the database.
func (r *userRepository) GetAllUsers(tx *gorm.DB) ([]User, error) {
	var users []User
	err := tx.Preload("Roles").Order("id ASC").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID retrieves a user by its ID from the database.
func (r *userRepository) GetUserByID(tx *gorm.DB, id int64) (User, error) {
	// Select the user with the given ID from the database
	var user User
	err := tx.Preload("Roles").First(&user, "id = ?", id).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errors.New("user with the given ID not found")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}

	return user, nil
}

// GetUserByUserName retrieves a user by their username from the database.
func (r *userRepository) GetUserByUserName(tx *gorm.DB, username string) (User, error) {
	// Select the user with the given username from the database
	var user User
	err := tx.Preload("Roles").First(&user, "lower(username) = lower(?)", username).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errors.New("user with the given username not found")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email from the database.
func (r *userRepository) GetUserByEmail(tx *gorm.DB, email string) (User, error) {
	// Select the user with the given email from the database
	var user User
	err := tx.Preload("Roles").First(&user, "lower(email) = lower(?)", email).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errors.New("user with the given email not found")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, err
	}

	return user, nil
}

// CreateUser inserts a new user into the database and returns the created user.
func (r *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, user User) (User, error) {
	// Insert the new user into the database
	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

// UpdateUser updates an existing user in the database and returns the updated user.
func (r *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user User) (User, error) {
	// Update the user in the database
	if err := tx.WithContext(ctx).Save(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}
