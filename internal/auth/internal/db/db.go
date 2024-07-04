// Package db provides functionality to connect to a SQLite database and Redis.
package db

import (
	"auth/ent"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Config holds the configuration for database connections.
type Config struct {
	// DatabaseURL is the connection string for the SQLite database.
	DatabaseURL string

	// RedisURL is the connection string for Redis.
	RedisURL string
}

// ConnectEnt establishes a connection to the SQLite database and returns an ent client.
// It also runs the auto migration tool to create the schema resources.
//
// Parameters:
//   - cfg: Config struct containing the DatabaseURL.
//
// Returns:
//   - *ent.Client: A pointer to the ent client for database operations.
//   - error: An error if the connection or schema creation fails, nil otherwise.
func ConnectEnt(cfg Config) (*ent.Client, error) {
	// Open a connection to the SQLite database
	client, err := ent.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool to create schema resources
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}

// ConnectRedis establishes a connection to Redis.
//
// Parameters:
//   - redisURL: The URL for connecting to Redis.
//
// Returns:
//   - *redis.Client: A pointer to the Redis client for Redis operations.
//   - error: An error if the connection fails, nil otherwise.
func ConnectRedis(redisURL string) (*redis.Client, error) {
	// Parse the Redis URL
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	// Create a new Redis client
	client := redis.NewClient(opts)

	// Ping Redis to check the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}