package refreshtoken

import (
	"context"

	"gorm.io/gorm"
)

// Interface for refresh token repository
// This interface defines the methods that the refresh token repository should implement
type RefreshTokenRepository interface {
	GetRefreshTokenByUserID(tx *gorm.DB, userID int64) (RefreshToken, error)
	GetRefreshTokenByToken(tx *gorm.DB, token string) (RefreshToken, error)
	CreateRefreshToken(ctx context.Context, tx *gorm.DB, token RefreshToken) (RefreshToken, error)
	RemoveRefreshTokenByUserID(ctx context.Context, tx *gorm.DB, userID int64) (bool, error)
}

// This struct defines the RefreshTokenRepository that contains methods for interacting with the database
// It implements the RefreshTokenRepository interface and provides methods for refresh token-related operations
type refreshTokenRepository struct{}

// NewRefreshTokenRepository creates a new instance of RefreshTokenRepository.
// It initializes the refreshTokenRepository struct and returns it.
func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{}
}

// GetRefreshTokenByUserID retrieves a refresh token by its user ID from the database.
func (r *refreshTokenRepository) GetRefreshTokenByUserID(tx *gorm.DB, userID int64) (RefreshToken, error) {
	// Select the refresh token with the given user ID from the database
	var refreshToken RefreshToken
	err := tx.First(&refreshToken, "user_id = ?", userID).Error
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

// GetRefreshTokenByToken retrieves a refresh token by its token string from the database.
func (r *refreshTokenRepository) GetRefreshTokenByToken(tx *gorm.DB, token string) (RefreshToken, error) {
	// Select the refresh token with the given token string from the database
	var refreshToken RefreshToken
	err := tx.First(&refreshToken, "token = ?", token).Error
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

// CreateRefreshToken creates a new refresh token in the database.
func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, tx *gorm.DB, token RefreshToken) (RefreshToken, error) {
	// Create a new refresh token in the database
	if err := tx.WithContext(ctx).Create(&token).Error; err != nil {
		return RefreshToken{}, err
	}

	return token, nil
}

// RemoveRefreshTokenByUserID removes a refresh token by its user ID from the database.
func (r *refreshTokenRepository) RemoveRefreshTokenByUserID(ctx context.Context, tx *gorm.DB, userID int64) (bool, error) {
	// Delete the refresh token with the given user ID from the database
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		return false, err
	}

	return true, nil
}
