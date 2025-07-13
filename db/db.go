package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Database configuration
const (
	defaultHost     = "postgres"
	defaultPort     = "5432"
	defaultUser     = "postgres"
	defaultDBName   = "todos"
	defaultPassword = "postgres"
	defaultSSLMode  = "disable"
	connectTimeout  = "20s"
)

// DB is the global database connection
var DB *sql.DB

// ConnectPostgres establishes a connection to PostgreSQL with retries
func ConnectPostgres() error {
	// Get connection parameters from environment variables
	host := getEnv("DB_HOST", defaultHost)
	port := getEnv("DB_PORT", defaultPort)
	user := getEnv("DB_USER", defaultUser)
	dbname := getEnv("DB_NAME", defaultDBName)
	password := getEnv("DB_PASSWORD", defaultPassword)
	sslmode := getEnv("DB_SSLMODE", defaultSSLMode)

	// Build connection string
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	// Try to connect with retries
	maxAttempts := 5
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempting to connect to PostgreSQL (attempt %d/%d)...\n", attempt, maxAttempts)

		// Open database connection
		var err error
		DB, err = sql.Open("postgres", connect)
		if err != nil {
			fmt.Printf("Attempt %d failed to open connection: %v\n", attempt, err)
			if attempt < maxAttempts {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := DB.PingContext(ctx); err != nil {
			fmt.Printf("Attempt %d failed to ping database: %v\n", attempt, err)
			if attempt < maxAttempts {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to ping database after %d attempts: %w", maxAttempts, err)
		}

		// If we get here, connection is successful
		fmt.Printf("Successfully connected to PostgreSQL on attempt %d\n", attempt)
		_, err = DB.Exec("CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY,title VARCHAR(255) NOT NULL,status BOOLEAN DEFAULT FALSE,deleted BOOLEAN DEFAULT FALSE);")
		if err != nil {
			return err
		}
		break
	}

	fmt.Println("Database connection and setup complete")
	return nil
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
