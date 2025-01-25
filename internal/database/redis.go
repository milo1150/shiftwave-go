package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Use the Redis service name in Docker Compose
		Password: password,
		DB:       0, // use default DB
	})
	return rdb
}

// Gen by chatgpt
// logs all keys and their values from Redis
func LogAllRedisKeysAndValues(ctx context.Context, rdb *redis.Client) error {
	var cursor uint64
	for {
		// Use SCAN to get a batch of keys
		keys, newCursor, err := rdb.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		// Iterate over the keys and log their values
		for _, key := range keys {
			value, err := rdb.Get(ctx, key).Result()
			if err == redis.Nil {
				// Key has expired or doesn't exist
				fmt.Printf("Key: %s no longer exists\n", key)
			} else if err != nil {
				return fmt.Errorf("failed to get key %s: %w", key, err)
			} else {
				fmt.Printf("Key: %s, Value: %s\n", key, value)
			}
		}

		// Update cursor for the next batch
		cursor = newCursor
		if cursor == 0 {
			break // No more keys to scan
		}
	}
	return nil
}

func GetRateLimitKey(ip, date string) string {
	return fmt.Sprintf("rate_limit:%s:%s", ip, date)
}

func GetJwtKey(username string) string {
	return fmt.Sprintf("jwt:%s", username)
}
