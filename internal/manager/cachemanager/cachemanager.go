package cachemanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"boilerplate-golang/internal/config"
)

var rdb *redis.Client
var ctx = context.Background()
// Init initializes the Redis client using config values.
// Redis connection is optional and will log a warning if it fails.
func Init() {
	cfg := config.Get()

	// Skip Redis initialization if host is not set
	if cfg.Redis.Host == "" {
		log.Println("Redis host not configured, skipping Redis connection")
		return
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  5 * time.Second, // Timeout for establishing a new connection
		ReadTimeout:  3 * time.Second, // Timeout for socket reads
		WriteTimeout: 3 * time.Second, // Timeout for socket writes
		PoolTimeout:  4 * time.Second, // Timeout for acquiring a connection from the pool
	})

	// Ping to ensure connection works; log warning if it fails
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Printf("Warning: Failed to ping Redis: %v", err)
		log.Println("Application will continue to run without Redis cache")
	}
	log.Println("cachemanager: Redis client initialized")
}

// Client returns the Redis client.
func Client() *redis.Client { return rdb }

// Set sets a value with expiration.
func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key.
func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}
