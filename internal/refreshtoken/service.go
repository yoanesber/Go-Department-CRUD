package refreshtoken

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"gorm.io/gorm"
)

var (
	JWTRefreshTokenExpirationHour string
)

// LoadEnv loads the environment variables.
func LoadEnv() {
	JWTRefreshTokenExpirationHour = os.Getenv("JWT_REFRESH_TOKEN_EXPIRATION_HOUR")
}

// This struct defines the RefreshTokenService that contains a repository field of type RefreshTokenRepository
// It implements the RefreshTokenService interface and provides methods for refresh token-related operations
type RefreshTokenService interface {
	GetRefreshTokenByUserID(ctx context.Context, userID int64) (RefreshToken, error)
	GetRefreshTokenByToken(ctx context.Context, token string) (RefreshToken, error)
	VerifyExpirationDate(ctx context.Context, exp time.Time) (bool, error)
	CreateRefreshToken(ctx context.Context, userID int64) (RefreshToken, error)
}

// This struct defines the RefreshTokenService that contains a repository field of type RefreshTokenRepository
// It implements the RefreshTokenService interface and provides methods for refresh token-related operations
type refreshTokenService struct {
	repo RefreshTokenRepository
}

// NewRefreshTokenService creates a new instance of RefreshTokenService with the given repository.
// It initializes the refreshTokenService struct and returns it.
func NewRefreshTokenService(repo RefreshTokenRepository) RefreshTokenService {
	return &refreshTokenService{repo: repo}
}

// GetRefreshTokenByUserID retrieves a refresh token by its user ID from the database.
func (s *refreshTokenService) GetRefreshTokenByUserID(ctx context.Context, userID int64) (RefreshToken, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return RefreshToken{}, errors.New("database connection is nil")
	}

	// Retrieve the token by user ID from the repository
	token, err := s.repo.GetRefreshTokenByUserID(db, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get refresh token by user ID: %v", err))
		return RefreshToken{}, err
	}

	return token, nil
}

// GetRefreshTokenByToken retrieves a refresh token by its token string from the database.
func (s *refreshTokenService) GetRefreshTokenByToken(ctx context.Context, token string) (RefreshToken, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return RefreshToken{}, errors.New("database connection is nil")
	}

	// Retrieve the token by token string from the repository
	refreshToken, err := s.repo.GetRefreshTokenByToken(db, token)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get refresh token by token: %v", err))
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

// VerifyExpirationDate checks if the expiration date is valid and not in the past.
func (s *refreshTokenService) VerifyExpirationDate(ctx context.Context, exp time.Time) (bool, error) {
	// Check if the expiration date is valid
	if exp.IsZero() {
		return false, errors.New("expiration date is zero")
	}

	// Check if the expiration date is in the past
	if time.Now().After(exp) {
		return false, nil
	}

	return true, nil
}

// CreateRefreshToken creates a new refresh token for the user in the database.
// If a refresh token already exists for the user, it will be removed before creating a new one,
// ensuring that only one refresh token exists for each user at a time.
func (s *refreshTokenService) CreateRefreshToken(ctx context.Context, userID int64) (RefreshToken, error) {
	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return RefreshToken{}, errors.New("database connection is nil")
	}

	var createdRefreshToken RefreshToken
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the refresh token already exists for the user
		existingRefreshToken, err := s.repo.GetRefreshTokenByUserID(tx, userID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// If the refresh token already exists, remove it
		if !existingRefreshToken.Equals(&RefreshToken{}) {
			if _, err := s.repo.RemoveRefreshTokenByUserID(ctx, tx, userID); err != nil {
				return err
			}
		}

		// Create a new refresh token
		tokenStr := uuid.New().String()
		refreshToken := RefreshToken{
			Token:      tokenStr,
			UserID:     userID,
			ExpiryDate: GetRefreshTokenExpiration(time.Now()),
		}

		// Create the refresh token in the database
		createdRefreshToken, err = s.repo.CreateRefreshToken(ctx, tx, refreshToken)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create refresh token: %v", err))
		return RefreshToken{}, err
	}

	return createdRefreshToken, nil
}

// GetRefreshTokenExpiration calculates the expiration date for the refresh token.
// It retrieves the expiration hour from an environment variable and adds it to the current time.
func GetRefreshTokenExpiration(now time.Time) time.Time {
	// Load environment variables
	LoadEnv()

	expHour, err := strconv.Atoi(JWTRefreshTokenExpirationHour)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse JWT_REFRESH_TOKEN_EXPIRATION_HOUR: %v", err))
		return now.Add(24 * time.Hour) // Default to 24 hours if the environment variable is not set or invalid
	}
	if expHour <= 0 {
		expHour = 24
	}

	return now.Add(time.Hour * time.Duration(expHour))
}
