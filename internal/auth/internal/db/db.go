package db

import (
	"auth/ent"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
}

// ConnectEnt establishes a connection to the SQLite database and returns an ent client
func ConnectEnt(cfg Config) (*ent.Client, error) {
	client, err := ent.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}

	// Run the auto migration tool
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}

	return client, nil
}

// ConnectRedis establishes a connection to Redis
func ConnectRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}
