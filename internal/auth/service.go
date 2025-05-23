package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yoanesber/Go-Department-CRUD/internal/refreshtoken"
	"github.com/yoanesber/Go-Department-CRUD/internal/role"
	"github.com/yoanesber/Go-Department-CRUD/internal/user"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/dbcontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util/redisutil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	JWTSecret         string
	TokenType         string
	SigningMethod     string
	JWTAudience       string
	JWTIssuer         string
	JWTExpirationHour string
	AccessTokenTTL    time.Duration
)

// LoadEnv loads environment variables
func LoadEnv() {
	JWTSecret = os.Getenv("JWT_SECRET")
	TokenType = os.Getenv("TOKEN_TYPE")
	SigningMethod = os.Getenv("JWT_ALGORITHM")
	JWTAudience = os.Getenv("JWT_AUDIENCE")
	JWTIssuer = os.Getenv("JWT_ISSUER")
	JWTExpirationHour = os.Getenv("JWT_EXPIRATION_HOUR")

	// Load access and refresh token TTL from environment variables
	access, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL_MINUTES"))
	AccessTokenTTL = time.Duration(access) * time.Minute
}

// Interface for auth service
// This interface defines the methods that the auth service should implement
type AuthService interface {
	Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error)
	RefreshToken(ctx context.Context, refreshTokenReq refreshtoken.RefreshTokenRequest) (refreshtoken.RefreshTokenResponse, error)
}

// This struct defines the AuthService that contains a user repository and a role repository
// It implements the AuthService interface and provides methods for authentication-related operations
type authService struct{}

// NewAuthService creates a new instance of AuthService with the given user and role repositories.
// It initializes the authService struct and returns it.
func NewAuthService() AuthService {
	return &authService{}
}

// Login authenticates a user with the given username and password.
// It retrieves the token for the user if the authentication is successful.
func (s *authService) Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error) {
	// Load environment variables
	LoadEnv()

	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return LoginResponse{}, errors.New("database connection is nil")
	}

	// Validate the authentication parameters using the validation
	if err := loginReq.Validate(); err != nil {
		return LoginResponse{}, err
	}

	var tokenStr string
	var refreshTokenStr string
	var expirationDateStr string
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the user exists
		userRepo := user.NewUserRepository()
		userService := user.NewUserService(userRepo)
		existingUser, err := userService.GetUserByUserName(ctx, loginReq.UserName)
		if err != nil {
			return err
		}

		// Check some conditions for the user
		if existingUser.Equals(&user.User{}) {
			return errors.New("user not found")
		}
		if !*existingUser.IsEnabled {
			return errors.New("user is not enabled")
		}
		if !*existingUser.IsAccountNonExpired {
			return errors.New("user account is expired")
		}
		if !*existingUser.IsAccountNonLocked {
			return errors.New("user account is locked")
		}
		if !*existingUser.IsCredentialsNonExpired {
			return errors.New("user credentials are expired")
		}
		if *existingUser.IsDeleted {
			return errors.New("user account is deleted")
		}

		// Compare the provided password with the stored hashed password
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginReq.Password)); err != nil {
			return errors.New("invalid password")
		}

		// Generate an access token for the user
		tokenStr, err = GenerateJWTToken(existingUser)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to generate JWT token: %v", err))
			return err
		}

		// Parse the JWT token
		jwtToken, err := ParseJWTToken(tokenStr)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to parse JWT token: %v", err))
			return err
		}

		// Get the expiration date from the token
		expirationDateStr, err = GetExpirationDateFromToken(jwtToken)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get expiration date from token: %v", err))
			return err
		}

		// Generate a refresh token for the user
		refreshTokenRepo := refreshtoken.NewRefreshTokenRepository()
		refreshTokenService := refreshtoken.NewRefreshTokenService(refreshTokenRepo)
		jwtRefreshToken, err := refreshTokenService.CreateRefreshToken(ctx, existingUser.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to create refresh token: %v", err))
			return err
		}
		if jwtRefreshToken.Equals(&refreshtoken.RefreshToken{}) {
			return errors.New("failed to create refresh token")
		}

		refreshTokenStr = jwtRefreshToken.Token

		// Update the last login time for the user
		_, err = userService.UpdateLastLogin(ctx, existingUser.ID, time.Now())
		if err != nil {
			logger.Error(fmt.Sprintf("failed to update last login time: %v", err))
			return err
		}

		// Store the access token details in Redis
		redisClient := dbcontext.GetRedisClient(ctx)
		if redisClient == nil {
			logger.Error("redis client is nil")
			return errors.New("redis client is nil")
		}
		redisKey := fmt.Sprintf("access_token:%s", existingUser.UserName)
		err = redisutil.SetJSON(ctx, redisClient, redisKey, LoginResponse{
			AccessToken:    tokenStr,
			RefreshToken:   refreshTokenStr,
			ExpirationDate: expirationDateStr,
			TokenType:      TokenType,
		}, AccessTokenTTL)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to set access token in Redis: %v", err))
			return err
		}

		return nil
	})

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		AccessToken:    tokenStr,
		RefreshToken:   refreshTokenStr,
		ExpirationDate: expirationDateStr,
		TokenType:      TokenType,
	}, nil
}

// RefreshToken refreshes the access token using the provided refresh token.
// It retrieves the new access token and refresh token for the user.
func (s *authService) RefreshToken(ctx context.Context, refreshTokenReq refreshtoken.RefreshTokenRequest) (refreshtoken.RefreshTokenResponse, error) {
	// Load environment variables
	LoadEnv()

	// Get the database connection from the context
	db := dbcontext.GetDB(ctx)
	if db == nil {
		logger.Error("database connection is nil")
		return refreshtoken.RefreshTokenResponse{}, errors.New("database connection is nil")
	}

	// Validate the refresh token request
	if err := refreshTokenReq.Validate(); err != nil {
		return refreshtoken.RefreshTokenResponse{}, err
	}

	var accessTokenStr string
	var refreshTokenStr string
	var expirationDateStr string
	err := db.Transaction(func(tx *gorm.DB) error {
		// Check if the refresh token exists
		refreshTokenRepo := refreshtoken.NewRefreshTokenRepository()
		refreshTokenService := refreshtoken.NewRefreshTokenService(refreshTokenRepo)
		existingRefreshToken, err := refreshTokenService.GetRefreshTokenByToken(ctx, refreshTokenReq.RefreshToken)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get refresh token: %v", err))
			return err
		}
		if existingRefreshToken.Equals(&refreshtoken.RefreshToken{}) {
			return errors.New("refresh token not found")
		}

		// If found, check if the refresh token is expired
		ok, _ := refreshTokenService.VerifyExpirationDate(ctx, existingRefreshToken.ExpiryDate)
		if !ok {
			return errors.New("refresh token is expired")
		}

		// Get user details using the user ID from the refresh token
		userRepo := user.NewUserRepository()
		userService := user.NewUserService(userRepo)
		userDetails, err := userService.GetUserByID(ctx, existingRefreshToken.UserID)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get user by ID: %v", err))
			return err
		}
		if userDetails.Equals(&user.User{}) {
			return errors.New("user not found")
		}

		// Generate an access token for the user
		accessTokenStr, err = GenerateJWTToken(userDetails)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to generate JWT token: %v", err))
			return err
		}

		// Parse the JWT token
		jwtToken, err := ParseJWTToken(accessTokenStr)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to parse JWT token: %v", err))
			return err
		}

		// Get the expiration date from the token
		expirationDateStr, err = GetExpirationDateFromToken(jwtToken)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get expiration date from token: %v", err))
			return err
		}

		// Regenerate a refresh token for the user
		jwtRefreshToken, err := refreshTokenService.CreateRefreshToken(ctx, userDetails.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to create refresh token: %v", err))
			return err
		}
		if jwtRefreshToken.Equals(&refreshtoken.RefreshToken{}) {
			return errors.New("failed to create refresh token")
		}

		refreshTokenStr = jwtRefreshToken.Token

		// Update the last login time for the user
		_, err = userService.UpdateLastLogin(ctx, userDetails.ID, time.Now())
		if err != nil {
			logger.Error(fmt.Sprintf("failed to update last login time: %v", err))
			return err
		}

		// Store the access token details in Redis
		redisClient := dbcontext.GetRedisClient(ctx)
		if redisClient == nil {
			logger.Error("redis client is nil")
			return errors.New("redis client is nil")
		}
		redisKey := fmt.Sprintf("access_token:%s", userDetails.UserName)
		err = redisutil.SetJSON(ctx, redisClient, redisKey, refreshtoken.RefreshTokenResponse{
			AccessToken:    accessTokenStr,
			RefreshToken:   refreshTokenStr,
			ExpirationDate: expirationDateStr,
			TokenType:      TokenType,
		}, AccessTokenTTL)

		if err != nil {
			logger.Error(fmt.Sprintf("failed to set access token in Redis: %v", err))
			return err
		}

		return nil
	})

	if err != nil {
		return refreshtoken.RefreshTokenResponse{}, err
	}

	return refreshtoken.RefreshTokenResponse{
		AccessToken:    accessTokenStr,
		RefreshToken:   refreshTokenStr,
		ExpirationDate: expirationDateStr,
		TokenType:      TokenType,
	}, nil
}

// GenerateJWTToken determines the function to use for generating a JWT token based on the signing method.
// It checks the signing method from the environment variable and calls the appropriate function.
func GenerateJWTToken(user user.User) (string, error) {
	// Load environment variables
	LoadEnv()

	// Check the signing method from the environment variable
	if SigningMethod == jwt.SigningMethodHS256.Alg() {
		return GenerateJWTTokenWithHS256(user)
	} else if SigningMethod == jwt.SigningMethodRS256.Alg() {
		return GenerateJWTTokenWithRS256(user)
	}

	return "", errors.New("unsupported signing method")
}

// GenerateJWTTokenWithHS256 generates a JWT token using the HS256 signing method.
// It creates the claims for the token and signs it with the secret key from the environment variable.
func GenerateJWTTokenWithHS256(user user.User) (string, error) {
	// Load environment variables
	LoadEnv()

	// Set the now time
	// This is used to set the issued at (iat) and expiration (exp) claims
	now := time.Now().Unix()

	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"sub":      user.UserName,
		"aud":      JWTAudience,
		"iss":      JWTIssuer,
		"iat":      now,
		"exp":      GetJWTExpiration(now),
		"email":    user.Email,
		"userid":   user.ID,
		"username": user.UserName,
		"roles":    ExtractRoleNames(user.Roles),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// GenerateJWTTokenWithRS256 generates a JWT token using the RS256 signing method.
// It creates the claims for the token and signs it with the private key loaded from the file.
func GenerateJWTTokenWithRS256(user user.User) (string, error) {
	// Load environment variables
	LoadEnv()

	// Load the private key from the file
	privateKey, err := util.LoadPrivateKey()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to load private key: %v", err))
		return "", err
	}

	// Set the now time
	// This is used to set the issued at (iat) and expiration (exp) claims
	now := time.Now().Unix()

	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"sub":      user.UserName,
		"aud":      JWTAudience,
		"iss":      JWTIssuer,
		"iat":      now,
		"exp":      GetJWTExpiration(now),
		"email":    user.Email,
		"userid":   user.ID,
		"username": user.UserName,
		"roles":    ExtractRoleNames(user.Roles),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ParseJWTToken determines the function to use for parsing a JWT token based on the signing method.
// It checks the signing method from the environment variable and calls the appropriate function.
func ParseJWTToken(tokenStr string) (*jwt.Token, error) {
	// Load environment variables
	LoadEnv()

	// Check the signing method from the environment variable
	if SigningMethod == jwt.SigningMethodHS256.Alg() {
		return ParseJWTTokenWithHS256(tokenStr)
	} else if SigningMethod == jwt.SigningMethodRS256.Alg() {
		return ParseJWTTokenWithRS256(tokenStr)
	}

	return nil, errors.New("unsupported signing method")
}

// ParseJWTTokenWithHS256 parses a JWT token using the HS256 signing method.
// It validates the token and returns the parsed token object.
func ParseJWTTokenWithHS256(tokenStr string) (*jwt.Token, error) {
	// Load environment variables
	LoadEnv()

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse JWT token: %v", err))
		return nil, err
	}
	return token, nil
}

// ParseJWTTokenWithRS256 parses a JWT token using the RS256 signing method.
// It validates the token and returns the parsed token object.
func ParseJWTTokenWithRS256(tokenStr string) (*jwt.Token, error) {
	// Load the public key from the file
	publicKey, err := util.LoadPublicKey()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to load public key: %v", err))
		return nil, err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse JWT token: %v", err))
		return nil, err
	}
	return token, nil
}

// GetRefreshTokenExpiration calculates the expiration time for the refresh token.
func GetJWTExpiration(now int64) int64 {
	// Load environment variables
	LoadEnv()

	expHour, err := strconv.Atoi(JWTExpirationHour)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse JWT expiration hour: %v", err))
		return now + int64(time.Hour.Seconds()*24)
	}
	if expHour <= 0 {
		expHour = 24
	}

	return now + int64(time.Duration(expHour)*time.Hour/time.Second)
}

// ExtractRoleNames extracts the role names from a slice of roles.
func ExtractRoleNames(roles []role.Role) []string {
	names := make([]string, len(roles))
	for i, r := range roles {
		names[i] = r.Name
	}
	return names
}

// GetExpirationDateFromToken extracts the expiration date from the JWT token claims.
func GetExpirationDateFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("failed to extract claims from token")
		return "", errors.New("failed to extract claims from token")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		logger.Error("failed to extract expiration date from claims")
		return "", errors.New("failed to extract expiration date from claims")
	}

	expirationDate := time.Unix(int64(expFloat), 0).Format(time.RFC3339)
	return expirationDate, nil
}
