// db/redis.go
package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var RCtx = context.Background()

func ConnectRedis() bool {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
	if err := RDB.Ping(RCtx).Err(); err != nil {
		return false
	}
	return true
}
