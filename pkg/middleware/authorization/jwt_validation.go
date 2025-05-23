package authorization

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yoanesber/Go-Department-CRUD/pkg/contextdata/metacontext"
	"github.com/yoanesber/Go-Department-CRUD/pkg/util"
)

var (
	TokenType string
	JWTSecret string
)

// LoadEnv loads environment variables
func LoadEnv() {
	TokenType = os.Getenv("TOKEN_TYPE")
	JWTSecret = os.Getenv("JWT_SECRET")
}

// JwtValidation is a middleware function that checks for a valid JWT token in the request header.
// It extracts the token from the "Authorization" header, validates it, and sets the user information in the context.
func JwtValidation() gin.HandlerFunc {
	// Load environment variables
	LoadEnv()

	return func(c *gin.Context) {
		// Get the token from the request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.JSONError(c, http.StatusUnauthorized, "No token provided", "Authorization header is missing")
			c.Abort()
			return
		}

		// Check if the token starts with TokenType
		tokenPrefix := TokenType + " "
		if !strings.HasPrefix(authHeader, tokenPrefix) {
			util.JSONError(c, http.StatusUnauthorized, "Invalid token format", fmt.Sprintf("Token must start with '%s'", tokenPrefix))
			c.Abort()
			return
		}

		// Extract the token string
		tokenStr := strings.TrimPrefix(authHeader, tokenPrefix)
		if tokenStr == "" {
			util.JSONError(c, http.StatusUnauthorized, "Invalid token format", "Token string is empty")
			c.Abort()
			return
		}

		// Parse the token and validate it
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// For HS256 signing method
			if token.Method.Alg() == jwt.SigningMethodHS256.Alg() {
				// Validate the token signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}

				// Return the secret key for validation
				return []byte(JWTSecret), nil
			}

			// For RS256 signing method
			// Load the public key from the environment variable
			publicKey, err := util.LoadPublicKey()
			if err != nil {
				return nil, err
			}

			// Validate the token signing method
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}

			// Return the public key for validation
			return publicKey, nil
		})

		if err != nil {
			util.JSONError(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		// Check if the token is valid
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			util.JSONError(c, http.StatusUnauthorized, "Invalid token", "Token is not valid")
			c.Abort()
			return
		}

		// Get the user ID from the claims
		// Convert the user ID to int64
		userID, _ := util.GetInt64Claim(claims, "userid")

		// Inject user information into the request context
		meta := metacontext.RequestMeta{
			UserID:   userID,
			UserName: claims["username"].(string),
			Email:    claims["email"].(string),
			Roles:    util.GetStringSliceClaim(claims, "roles"),
		}
		ctx := metacontext.InjectRequestMeta(c.Request.Context(), meta)

		// Set the new request context with user information
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
