package storage

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")

	// Debug (you can remove later)
	fmt.Println("REDIS_URL:", redisURL)

	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic("invalid REDIS_URL: " + err.Error())
		}
		return redis.NewClient(opt)
	}

	// fallback for local only
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	fmt.Println("Using REDIS_ADDR:", addr)

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}