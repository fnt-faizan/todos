package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectPostgres connects to PostgreSQL with a 30-second timeout
func ConnectPostgres() error {
	// Get connection parameters from environment variables with defaults
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	dbname := getEnv("DB_NAME", "todos")
	password := getEnv("DB_PASSWORD", "")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Set connection timeout
	connectTimeout := "30s"

	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s connect_timeout=%s",
		host, port, user, password, dbname, sslmode, connectTimeout)

	var err error
	DB, err = sql.Open("postgres", connect)
	if err != nil {
		return fmt.Errorf("error opening database %s", err)
	}

	// Set connection pool settings
	DB.SetConnMaxLifetime(30 * time.Second)
	DB.SetConnMaxIdleTime(30 * time.Second)
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(25)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to connect to database %s", err)
	}

	return nil
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
