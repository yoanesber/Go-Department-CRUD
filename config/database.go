package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-gorp/gorp" // Gorp is a Go library for working with SQL databases
	_ "github.com/lib/pq"     // PostgreSQL driver
	// Importing the model package for the Department struct
)

// DB struct represents the database connection
// It embeds the sql.DB type to provide database functionalities
type DB struct {
	*sql.DB
}

// gorp is a package that provides a simple ORM for Go
var db *gorp.DbMap

// InitDB initializes the database connection using environment variables
// It constructs the connection string and calls ConnectDB to establish the connection
func InitDB() {
	// Create the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_SSL"))

	// Initialize the database connection
	var err error
	db, err = ConnectDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// ConnectDB initializes the database connection
func ConnectDB(connStr string) (*gorp.DbMap, error) {
	// Open a new database connection
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	//dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests

	return dbMap, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if db != nil {
		if err := db.Db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}
}

// GetDB returns the database connection
func GetDB() *gorp.DbMap {
	return db
}
