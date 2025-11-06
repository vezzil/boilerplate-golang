package redismanager

import (
	"context"
	"fmt"
	"log"
	"time"

	"boilerplate-golang/internal/infrastructure/config"

	"github.com/redis/go-redis/v9"
)

var Redis *RedisClient

// RedisClient wraps the redis.Client
type RedisClient struct {
	*redis.Client
}

func Init() {
	cfg := config.Get()

	// Skip Redis initialization if host is not set
	if cfg.Redis.Host == "" {
		log.Println("Redis host not configured, skipping Redis connection")
		return
	}

	Redis = &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			DialTimeout:  5 * time.Second, // Timeout for establishing a new connection
			ReadTimeout:  3 * time.Second, // Timeout for socket reads
			WriteTimeout: 3 * time.Second, // Timeout for socket writes
			PoolTimeout:  4 * time.Second, // Timeout for acquiring a connection from the pool
		}),
	}

	// Ping to ensure connection works; log warning if it fails
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := Redis.Ping(ctx).Result(); err != nil {
		log.Printf("Warning: Failed to ping Redis: %v", err)
		log.Println("Application will continue to run without Redis cache")
	} else {
		log.Println("redismanager: Redis client initialized")
	}
}

// RSet sets a value with expiration in seconds.
func (r *RedisClient) RSet(key string, value interface{}, ex int) error {
	return r.Set(context.Background(), key, value, time.Duration(ex)*time.Second).Err()
}

// RGet retrieves a value by key.
func (r *RedisClient) RGet(key string) string {
	value, err := r.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return value
}

// RTTL returns the expiration time of the key in seconds.
func (r *RedisClient) RTTL(key string) int {
	value, err := r.TTL(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	return int(value.Seconds())
}

// RDel deletes a key.
func (r *RedisClient) RDel(key string) error {
	return r.Del(context.Background(), key).Err()
}

// Close closes the Redis client.
func (r *RedisClient) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}

// GetRedisClient gets or initializes the Redis client.
func GetRedisClient() (*RedisClient, error) {
	if Redis == nil {
		Init()
		if Redis == nil {
			return nil, fmt.Errorf("failed to initialize Redis client")
		}
	}
	return Redis, nil
}
