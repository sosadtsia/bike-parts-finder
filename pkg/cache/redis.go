package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sosadtsia/bike-parts-finder/pkg/models"
)

// RedisClient is a wrapper around Redis client
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new RedisClient
func NewRedisClient() (*RedisClient, error) {
	// Get Redis connection string from environment variable
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parsing Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Check connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connecting to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
	}, nil
}

// Close closes the Redis client
func (c *RedisClient) Close() error {
	return c.client.Close()
}

// Get retrieves a value by key
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set stores a value with a key and expiration
func (c *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// CachePart caches a part
func (c *RedisClient) CachePart(ctx context.Context, part models.Part) error {
	key := fmt.Sprintf("part:%s", part.ID)
	data, err := json.Marshal(part)
	if err != nil {
		return fmt.Errorf("marshaling part %s: %w", part.ID, err)
	}

	return c.Set(ctx, key, string(data), 24*time.Hour)
}

// GetCachedPart retrieves a cached part
func (c *RedisClient) GetCachedPart(ctx context.Context, id string) (models.Part, error) {
	var part models.Part
	key := fmt.Sprintf("part:%s", id)

	data, err := c.Get(ctx, key)
	if err != nil {
		return part, fmt.Errorf("getting part %s from cache: %w", id, err)
	}

	if err := json.Unmarshal([]byte(data), &part); err != nil {
		return part, fmt.Errorf("unmarshaling part %s: %w", id, err)
	}

	return part, nil
}
