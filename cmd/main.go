package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/yoanesber/Go-Department-CRUD/config/db/postgresdb"
	"github.com/yoanesber/Go-Department-CRUD/config/db/redisdb"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"github.com/yoanesber/Go-Department-CRUD/pkg/validator"
	"github.com/yoanesber/Go-Department-CRUD/routes"
)

var (
	Environment string
	Port        string
	IsSSL       string
	APIVersion  string
	SSLKeys     string
	SSLCert     string
)

// Init function to initialize the application
// This function is called when the application starts
func init() {
	// Init logger
	logger.InitLoggers()
}

// Main function to start the Gin server and set up routes
// It loads environment variables, sets up middleware, and defines API routes
func main() {
	// Load environment variables from .env file
	// _ = godotenv.Load(".env")

	// Get environment variable from .env file
	Environment := os.Getenv("ENV")
	Port := os.Getenv("PORT")
	IsSSL := os.Getenv("IS_SSL")
	APIVersion := os.Getenv("API_VERSION")
	SSLKeys := os.Getenv("SSL_KEYS")
	SSLCert := os.Getenv("SSL_CERT")

	// Set the Gin mode based on the environment
	gin.SetMode(gin.DebugMode)
	if Environment == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the PostgreSQL database connection using the configuration from the .env file
	postgresdb.LoadEnv()
	postgresdb.InitDB()

	// Initialize the Redis client using the configuration from the .env file
	redisdb.LoadEnv()
	redisdb.InitRedis()

	// Initialize the validator for request validation
	validator.InitValidator()

	// Set up Gin server with middleware and routes
	r := routes.SetupRouter()

	// Set up trusted proxies for Gin
	// This is used to trust the X-Forwarded-For header for client IP detection
	r.SetTrustedProxies(nil)

	if Port == "" {
		Port = "8080" // Default port if not specified in .env
	}

	// Log the server start information
	logger.Info("Starting server on : ", log.Fields{
		"port":    Port,
		"env":     Environment,
		"ssl":     IsSSL,
		"version": APIVersion,
	})

	// Start the server with or without SSL based on the environment variable
	if IsSSL == "TRUE" {
		//Generated using sh generate-certificate.sh
		err := r.RunTLS(":"+Port, SSLCert, SSLKeys)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start server with SSL: %v", err))
		}
	} else {
		err := r.Run(":" + Port)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		}
	}
}
