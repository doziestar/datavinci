// Package db provides functionality to connect to a SQLite database and Redis.
package db

import (
	"auth/ent"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
)

// Config holds the configuration for database connections.
type Config struct {
	// DatabaseURL is the connection string for the SQLite database.
	DatabaseURL string

	// RedisURL is the connection string for Redis.
	RedisURL string

	// MigrationsDir is the directory containing the ent migration files.
	MigrationsDir string
}

// ConnectEnt establishes a connection to the SQLite database and returns an ent client.
// It checks the ent folder and applies the latest migration.
//
// Parameters:
//   - cfg: Config struct containing the DatabaseURL and MigrationsDir.
//
// Returns:
//   - *ent.Client: A pointer to the ent client for database operations.
//   - error: An error if the connection or migration fails, nil otherwise.
func ConnectEnt(cfg Config) (*ent.Client, error) {
	// Open a connection to the SQLite database
	client, err := ent.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	// Find the latest migration file
	latestMigration, err := getLatestMigrationFile(cfg.MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest migration file: %v", err)
	}

	log.Println("latestMigration", latestMigration)

	// // Apply the latest migration
	// err = client.Schema.Create(context.Background(), ent.Migrate(latestMigration))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to apply latest migration: %v", err)
	// }

	return client, nil
}

// getLatestMigrationFile finds the most recent migration file in the specified directory.
//
// Parameters:
//   - dir: The directory containing migration files.
//
// Returns:
//   - string: The content of the latest migration file.
//   - error: An error if reading the migration file fails, nil otherwise.
func getLatestMigrationFile(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read migration directory: %v", err)
	}

	var latestFile os.DirEntry
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			if latestFile == nil || file.Name() > latestFile.Name() {
				latestFile = file
			}
		}
	}

	if latestFile == nil {
		return "", fmt.Errorf("no migration files found in directory")
	}

	content, err := os.ReadFile(filepath.Join(dir, latestFile.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to read migration file: %v", err)
	}

	return string(content), nil
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