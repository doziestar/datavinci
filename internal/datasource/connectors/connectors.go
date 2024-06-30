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

	// Query executes a query and returns the results.
	Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error)

	// Execute executes a command (e.g., INSERT, UPDATE, DELETE) and returns the number of affected rows.
	Execute(ctx context.Context, command string, args ...interface{}) (int64, error)

	// Ping checks if the data source is accessible.
	Ping(ctx context.Context) error
}

// Config represents the configuration for a data connector.
type Config struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Options  map[string]interface{}
}