package cache

import (
	"context"
	"fmt"
	"time"

	"go-modular-monolith/pkg/framework"
	"go-modular-monolith/pkg/framework/providers"

	"github.com/redis/go-redis/v9"
)

// RedisProvider provides Redis caching functionality
type RedisProvider struct {
	client *redis.Client
	config *framework.FrameworkConfig
}

// NewRedisProvider creates a new Redis provider
func NewRedisProvider() *RedisProvider {
	return &RedisProvider{}
}

// Name returns the provider name
func (r *RedisProvider) Name() string {
	return "redis"
}

// IsEnabled checks if Redis is enabled in configuration
func (r *RedisProvider) IsEnabled() bool {
	return framework.IsEnabled("cache", "redis")
}

// Initialize initializes the Redis connection
func (r *RedisProvider) Initialize(ctx context.Context, deps providers.Dependencies) error {
	r.config = deps.Config

	// Get Redis configuration from environment or config
	redisURL := getRedisURL()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	r.client = redis.NewClient(opts)

	// Test connection
	return r.HealthCheck(ctx)
}

// Shutdown closes the Redis connection
func (r *RedisProvider) Shutdown(ctx context.Context) error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// HealthCheck pings Redis to check connectivity
func (r *RedisProvider) HealthCheck(ctx context.Context) error {
	if r.client == nil {
		return fmt.Errorf("Redis client not initialized")
	}

	pong := r.client.Ping(ctx)
	if pong.Err() != nil {
		return fmt.Errorf("Redis health check failed: %v", pong.Err())
	}

	return nil
}

// GetClient returns the Redis client instance
func (r *RedisProvider) GetClient() *redis.Client {
	return r.client
}

// Cache interface implementation
func (r *RedisProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisProvider) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisProvider) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisProvider) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// getRedisURL gets Redis connection URL from environment or returns default
func getRedisURL() string {
	// This would typically read from environment variables or config
	// For now, returning a default local Redis URL
	return "redis://localhost:6379/0"
}
