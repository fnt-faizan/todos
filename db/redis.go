// db/redis.go
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// ConnectRedis connects to Redis with a 30-second timeout
func ConnectRedis() error {
	// Get Redis connection parameters from environment variables with defaults
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")

	addr := fmt.Sprintf("%s:%s", host, port)

	// Create Redis client with timeout settings
	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DialTimeout:  30 * time.Second, // Timeout for establishing connection
		ReadTimeout:  30 * time.Second, // Timeout for reading data
		WriteTimeout: 30 * time.Second, // Timeout for writing data
		PoolSize:     10,               // Maximum number of connections
		MinIdleConns: 2,                // Minimum number of idle connections
	})

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}

	return nil
}
