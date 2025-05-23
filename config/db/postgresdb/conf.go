package postgresdb

import (
	"fmt"
	"os"

	"github.com/yoanesber/Go-Department-CRUD/internal/department"
	"github.com/yoanesber/Go-Department-CRUD/internal/refreshtoken"
	"github.com/yoanesber/Go-Department-CRUD/internal/role"
	"github.com/yoanesber/Go-Department-CRUD/internal/user"
	"github.com/yoanesber/Go-Department-CRUD/pkg/logger"
	"gorm.io/driver/postgres"        // Import the PostgreSQL driver for GORM
	"gorm.io/gorm"                   // Import GORM for ORM functionalities
	gormLogger "gorm.io/gorm/logger" // Import GORM logger for logging SQL queries
)

var (
	db         *gorm.DB
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBSchema   string
	DBSSL      string
	DBTimeZone string
	DBMigrate  string
	DBSeed     string
	DBSeedFile string
	DBLog      string
)

// LoadEnv loads environment variables from the .env file
// It sets the database connection parameters such as host, port, user, password, etc.
func LoadEnv() {
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBUser = os.Getenv("DB_USER")
	DBPass = os.Getenv("DB_PASS")
	DBName = os.Getenv("DB_NAME")
	DBSchema = os.Getenv("DB_SCHEMA")
	DBSSL = os.Getenv("DB_SSL")
	DBTimeZone = os.Getenv("DB_TIMEZONE")
	DBMigrate = os.Getenv("DB_MIGRATE")
	DBSeed = os.Getenv("DB_SEED")
	DBSeedFile = os.Getenv("DB_SEED_FILE")
	DBLog = os.Getenv("DB_LOG")
}

// InitDB initializes the GORM database connection
func InitDB() {
	// Create the connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		DBHost,
		DBPort,
		DBUser,
		DBPass,
		DBName,
		DBSSL,
		DBTimeZone,
	)

	// Set the log level based on the environment variable
	var logLevel gormLogger.LogLevel
	if DBLog == "INFO" {
		logLevel = gormLogger.Info
	} else if DBLog == "ERROR" {
		logLevel = gormLogger.Error
	} else if DBLog == "SILENT" {
		logLevel = gormLogger.Silent
	} else {
		logLevel = gormLogger.Warn
	}

	// Open the connection using GORM and PostgreSQL driver
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
		return
	}

	logger.Info("Connected to PostgreSQL database")

	// Migrate the database schema
	if DBMigrate == "TRUE" {
		err := db.Transaction(func(tx *gorm.DB) error {
			// Drop and recreate tables if they exist
			err = tx.Migrator().DropTable(&refreshtoken.RefreshToken{}, &role.UserRole{}, &role.Role{}, &user.User{}, &department.Department{})
			if err != nil {
				return fmt.Errorf("failed to drop tables: %v", err)
			}

			// Migrate the database schema
			err = tx.AutoMigrate(&role.Role{}, &user.User{}, &refreshtoken.RefreshToken{}, &department.Department{})
			if err != nil {
				return fmt.Errorf("failed to migrate database: %v", err)
			}

			if DBSeed == "TRUE" {
				// Import initial data from the seed file
				if DBSeedFile == "" {
					return fmt.Errorf("DB_SEED_FILE environment variable is not set")
				}

				// Read the seed file
				seedData, err := os.ReadFile(DBSeedFile)
				if err != nil {
					return fmt.Errorf("failed to read seed file: %v", err)
				}

				// Execute the seed data
				if err := tx.Exec(string(seedData)).Error; err != nil {
					return fmt.Errorf("failed to execute seed data: %v", err)
				}
			}

			return nil
		})

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to migrate database: %v", err))
			return
		}

		logger.Info("Database migrated successfully")
	}
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return db
}

// CloseDB closes the database connection (optional, for when needed)
func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get SQL DB: %v", err))
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Error(fmt.Sprintf("Failed to close database connection: %v", err))
	}
}
