package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"github.com/joho/godotenv" // gin-swagger middleware
	"github.com/yoanesber/Go-Department-CRUD/config"
	"github.com/yoanesber/Go-Department-CRUD/routes"
)

// Middleware to handle Cross-Origin Resource Sharing (CORS) for the API
// This middleware allows requests from specific origins and handles preflight requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}

// Middleware to generate a unique request ID for each incoming request
// This ID is added to the response headers for tracking purposes
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}

// Main function to start the Gin server and set up routes
// It loads environment variables, sets up middleware, and defines API routes
func main() {
	// Load environment variables from .env file
	_ = godotenv.Load(".env")

	// Get ENV variable from .env file
	env := os.Getenv("ENV")
	if env == "" {
		env = "DEVELOPMENT" // Default to DEVELOPMENT if ENV is not set
	}

	// Determine the environment file to load based on the ENV variable
	// This allows for different configurations for different environments (e.g., development, staging, production)
	var envFile string
	if env == "DEVELOPMENT" {
		envFile = ".env.development" // Use development environment file
	} else if env == "STAGING" {
		envFile = ".env.staging" // Use staging environment file
	} else {
		envFile = ".env.production" // Use default environment file
	}

	// Load the specified environment file if it exists
	if envFile != "" {
		err := godotenv.Overload(envFile)
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	// Set the Gin mode based on the environment
	if env == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the database connection using the configuration from the .env file
	config.InitDB()

	// Set up Gin server with middleware and routes
	r := routes.SetupRouter()
	r.Use(
		CORSMiddleware(),
		RequestIDMiddleware(),
		gzip.Gzip(gzip.DefaultCompression),
	)

	// Set up trusted proxies for Gin
	// This is used to trust the X-Forwarded-For header for client IP detection
	r.SetTrustedProxies(nil)

	// Set up logging for the server
	// This middleware logs the details of each request and response
	// It includes information such as client IP, request method, path, user agent, status code, latency, and error message (if any)
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - %s - %s - %s - %d - %s - %s\n",
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.UserAgent(),
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {
		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/mycert.cer",
			KEY:  "./cert/mycert.key",
		}

		err := r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
		if err != nil {
			log.Fatalf("Failed to start server with SSL: %v", err)
		}
	} else {
		err := r.Run(":" + port)
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
