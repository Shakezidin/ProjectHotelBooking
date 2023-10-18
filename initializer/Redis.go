package initializer

import (
	"os"

	"github.com/redis/go-redis/v9"
)

// ReddisClient initializes and returns a Redis client.
var ReddisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("RedisAddr"),
	Password: os.Getenv("RedisPass"),
	DB:       0,
})
