package util

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/go-playground/validator.v9"
)

var (
	JWTPublicKeyPath  string
	JWTPrivateKeyPath string
)

// LoadEnv loads environment variables
func LoadEnv() {
	JWTPublicKeyPath = os.Getenv("JWT_PUBLIC_KEY_PATH")
	JWTPrivateKeyPath = os.Getenv("JWT_PRIVATE_KEY_PATH")
}

// FormatValidationErrors formats validation errors into a slice of maps.
// Each map contains the field name and the corresponding error message.
func FormatValidationErrors(err error) []map[string]string {
	var errors []map[string]string

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			// Customize the message based on tag
			var message string
			switch fe.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", fe.Field())
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", fe.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
			default:
				message = fmt.Sprintf("%s is not valid", fe.Field())
			}

			errors = append(errors, map[string]string{
				"field":   fe.Field(),
				"message": message,
			})
		}
	}
	return errors
}

// LoadPublicKey loads the public key from the specified path in the environment variable.
// It returns the parsed RSA public key or an error if the file cannot be read or parsed.
func LoadPublicKey() (*rsa.PublicKey, error) {
	// Load environment variables
	LoadEnv()

	keyData, err := os.ReadFile(JWTPublicKeyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

// LoadPrivateKey loads the private key from the specified path in the environment variable.
// It returns the parsed RSA private key or an error if the file cannot be read or parsed.
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	// Load environment variables
	LoadEnv()

	keyData, err := os.ReadFile(JWTPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}

// GetInt64Claim retrieves an int64 claim from the JWT claims.
// It checks if the claim exists and is of type float64, then converts it to int64.
func GetInt64Claim(claims jwt.MapClaims, key string) (int64, error) {
	if val, ok := claims[key]; ok {
		if f, ok := val.(float64); ok {
			return int64(f), nil
		}
		return 0, fmt.Errorf("claim %s is not a float64", key)
	}
	return 0, fmt.Errorf("claim %s not found", key)
}

// GetStringClaim retrieves a string claim from the JWT claims.
// It checks if the claim exists and is of type string.
func GetStringSliceClaim(claims jwt.MapClaims, key string) []string {
	if val, ok := claims[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			strSlice := make([]string, len(slice))
			for i, v := range slice {
				if str, ok := v.(string); ok {
					strSlice[i] = str
				}
			}
			return strSlice
		}
	}
	return nil
}
