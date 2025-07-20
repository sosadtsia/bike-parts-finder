package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/svosadtsia/bike-parts-finder/pkg/models"
)

// RedisConfig holds the configuration for connecting to Redis
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Cache represents a cache for bike parts search results
type Cache struct {
	client *redis.Client
}

// NewRedisCache creates a new cache with a Redis connection
func NewRedisCache(config RedisConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Ping the Redis server to check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// cacheKey generates a cache key for search parameters
func cacheKey(params models.SearchParams) string {
	return fmt.Sprintf("search:%s:%s:%d:%s:%d:%d",
		params.Brand, params.Model, params.Year, params.Category, params.Page, params.Limit)
}

// GetParts tries to retrieve bike parts from the cache
func (c *Cache) GetParts(ctx context.Context, params models.SearchParams) ([]models.Part, int, bool, error) {
	key := cacheKey(params)

	// Try to get cached data
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Cache miss, not an error
			return nil, 0, false, nil
		}
		return nil, 0, false, fmt.Errorf("redis get error: %w", err)
	}

	// Unmarshal the cached data
	var cachedResult struct {
		Parts []models.Part `json:"parts"`
		Total int           `json:"total"`
	}

	if err := json.Unmarshal(data, &cachedResult); err != nil {
		return nil, 0, false, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return cachedResult.Parts, cachedResult.Total, true, nil
}

// SetParts caches bike parts search results
func (c *Cache) SetParts(ctx context.Context, params models.SearchParams, parts []models.Part, total int) error {
	key := cacheKey(params)

	// Create a struct with both parts and total count
	cacheData := struct {
		Parts []models.Part `json:"parts"`
		Total int           `json:"total"`
	}{
		Parts: parts,
		Total: total,
	}

	// Marshal the data
	data, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("failed to marshal data for cache: %w", err)
	}

	// Cache for 15 minutes
	expiration := 15 * time.Minute

	// Set in Redis
	if err := c.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to cache data: %w", err)
	}

	return nil
}

// InvalidateCache removes all search-related cache entries
func (c *Cache) InvalidateCache(ctx context.Context) error {
	// Get all keys matching the search pattern
	keys, err := c.client.Keys(ctx, "search:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get cache keys: %w", err)
	}

	// No keys to delete
	if len(keys) == 0 {
		return nil
	}

	// Delete all matching keys
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	return nil
}
