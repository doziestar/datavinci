package connectors

import (
	"context"
	"fmt"

	"pkg/common/errors"
	"pkg/common/retry"

	"github.com/go-redis/redis/v8"
)

// RedisConnector implements the Connector interface for Redis.
type RedisConnector struct {
	client *redis.Client
	config *Config
}

// NewRedisConnector creates a new RedisConnector with the given configuration.
func NewRedisConnector(config *Config) *RedisConnector {
	return &RedisConnector{config: config}
}

// Connect establishes a connection to the Redis database.
func (c *RedisConnector) Connect(ctx context.Context) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
		Password: c.config.Password,
		DB:       c.config.RedisDB, // Redis uses integer for database selection
	})

	err := retry.Retry(ctx, func() error {
		return client.Ping(ctx).Err()
	}, retry.DefaultConfig())

	if err != nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, "failed to connect to Redis", err)
	}

	c.client = client
	return nil
}

// Close closes the connection to the Redis database.
func (c *RedisConnector) Close(ctx context.Context) error {
	if c.client == nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, "connection already closed", nil)
	}
	return c.client.Close()
}

// Query executes a query and returns the results as a slice of maps.
// For Redis, this is implemented as a key lookup.
func (c *RedisConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if c.client == nil {
		return nil, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	val, err := c.client.Get(ctx, query).Result()
	if err == redis.Nil {
		return nil, nil // Key does not exist
	} else if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to execute query", err)
	}

	return []map[string]interface{}{{"value": val}}, nil
}

// Execute executes a command and returns the number of affected keys.
func (c *RedisConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	if c.client == nil {
		return 0, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	cmd := c.client.Do(ctx, append([]interface{}{command}, args...)...)
	if cmd.Err() != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to execute command", cmd.Err())
	}

	// For simplicity, we're returning 1 if the command was successful.
	// In a real-world scenario, you might want to interpret the result based on the specific command.
	return 1, nil
}

// Ping checks if the database connection is still alive.
func (c *RedisConnector) Ping(ctx context.Context) error {
	if c.client == nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}
	return c.client.Ping(ctx).Err()
}

// Transaction starts a new database transaction and returns a TransactionConnector.
func (c *RedisConnector) Transaction(ctx context.Context) (TransactionConnector, error) {
	if c.client == nil {
		return nil, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	tx := c.client.TxPipeline()
	return &RedisTransactionConnector{tx: tx}, nil
}

// RedisTransactionConnector implements the TransactionConnector interface for Redis.
type RedisTransactionConnector struct {
	tx redis.Pipeliner
}

// Query executes a query within the transaction and returns the results.
func (c *RedisTransactionConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	cmd := c.tx.Get(ctx, query)
	val, err := cmd.Result()
	if err == redis.Nil {
		return nil, nil // Key does not exist
	} else if err != nil {
		return nil, errors.NewError(errors.ErrorTypeTransaction, "cannot query in Redis transaction, use Execute instead", nil)
	}

	return []map[string]interface{}{{"value": val}}, nil
}

// Execute executes a command within the transaction.
func (c *RedisTransactionConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	cmd := c.tx.Do(ctx, append([]interface{}{command}, args...)...)
	if cmd.Err() != nil {
		return 0, errors.NewError(errors.ErrorTypeTransaction, "failed to execute command", cmd.Err())
	}
	return 0, nil // The actual execution happens on Commit
}

// Commit commits the transaction.
func (c *RedisTransactionConnector) Commit(ctx context.Context) error {
	_, err := c.tx.Exec(ctx)
	if err != nil {
		return errors.NewError(errors.ErrorTypeTransaction, "failed to commit transaction", err)
	}
	return nil
}

// Rollback rolls back the transaction.
func (c *RedisTransactionConnector) Rollback(ctx context.Context) error {
	err := c.tx.Discard()
	if err != nil {
		return errors.NewError(errors.ErrorTypeTransaction, "failed to rollback transaction", err)
	}
	return nil
}
