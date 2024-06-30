package connectors

import (
	"context"
)

// Connector is the interface that wraps the basic methods for a data connector.
type Connector interface {
	// Connect establishes a connection to the data source.
	Connect(ctx context.Context) error

	// Close closes the connection to the data source.
	Close(ctx context.Context) error

	// Query executes a query that returns rows, typically a SELECT statement.
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)

	// Execute executes a command that does not return rows, typically an INSERT, UPDATE, DELETE, or CREATE statement.
	Execute(ctx context.Context, command string, args ...interface{}) (int64, error)

	// Ping verifies a connection to the data source.
	Ping(ctx context.Context) error

	// BeginTransaction starts a new transaction.
	Transaction(ctx context.Context) (TransactionConnector, error)
}

// TransactionConnector is the interface for transaction-specific operations.
type TransactionConnector interface {
	// Query executes a query that returns rows, typically a SELECT statement.
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)

	// Execute executes a command that does not return rows, typically an INSERT, UPDATE, DELETE, or CREATE statement.
	Execute(ctx context.Context, command string, args ...interface{}) (int64, error)

	// Commit commits the transaction.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction.
	Rollback(ctx context.Context) error
}

// Config represents the configuration for a data connector.
type Config struct {
	// Type is the type of the data connector.
	Type                    string
	// Host is the hostname of the data source.
	Host                    string
	// Port is the port number of the data source.
	Port                    int
	// Username is the username for the data source.
	Username                string
	// Password is the password for the data source.
	Password                string
	// Database is the name of the database.
	Database                string
	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns            int
	// MaxIdleConns is the maximum number of connections in the idle connection pool.
	MaxIdleConns            int
	// ConnMaxLifetimeSeconds is the maximum amount of time a connection may be reused.
	ConnMaxLifetimeSeconds  int
	// ConnMaxIdleTimeSeconds is the maximum amount of time a connection may be idle before being closed.
	Options                 map[string]interface{}
}