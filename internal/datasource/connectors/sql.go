// Package connectors provides implementations of various data source connectors.
package connectors

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pkg/common/errors"
	"pkg/common/retry"

	_ "github.com/lib/pq"
)

// SQLConnector implements the Connector interface for SQL databases.
type SQLConnector struct {
	db     *sql.DB
	config *Config
}

// NewSQLConnector creates a new SQLConnector with the given configuration.
func NewSQLConnector(config *Config) *SQLConnector {
	return &SQLConnector{config: config}
}

// Connect establishes a connection to the SQL database.
func (c *SQLConnector) Connect(ctx context.Context) error {
	var dsn string
	var driverName string

	switch c.config.Driver {
	case "postgres":
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.config.Host, c.config.Port, c.config.Username, c.config.Password, c.config.Database)
	case "mysql":
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			c.config.Username, c.config.Password, c.config.Host, c.config.Port, c.config.Database)
	case "sqlite":
		driverName = "sqlite3"
		dsn = c.config.Database // For SQLite, the database is the file path
	default:
		return errors.NewError(errors.ErrorTypeConfiguration, "unsupported SQL driver", nil)
	}

	var db *sql.DB
	err := retry.Retry(ctx, func() error {
		var err error
		db, err = sql.Open(driverName, dsn)
		if err != nil {
			return errors.NewError(errors.ErrorTypeDatabaseConnection, "failed to open database", err)
		}
		return db.PingContext(ctx)
	}, retry.DefaultConfig())

	if err != nil {
		return err
	}

	// Configure connection pool
	db.SetMaxOpenConns(c.config.MaxOpenConns)
	db.SetMaxIdleConns(c.config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(c.config.ConnMaxLifetimeSeconds) * time.Second)

	c.db = db
	return nil
}

// Close closes the connection to the SQL database.
func (c *SQLConnector) Close(ctx context.Context) error {
	if c.db == nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, "connection already closed", nil)
	}
	return c.db.Close()
}

// Query executes a query and returns the results as a slice of maps.
func (c *SQLConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if c.db == nil {
		return nil, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to execute query", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to get columns", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to scan row", err)
		}

		row := make(map[string]interface{})
		for i, column := range columns {
			row[column] = *(values[i].(*interface{}))
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "error during row iteration", err)
	}

	return results, nil
}

// Execute executes a command (e.g., INSERT, UPDATE, DELETE) and returns the number of affected rows.
func (c *SQLConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	if c.db == nil {
		return 0, errors.NewError(errors.ErrorTypeDatabaseConnection,errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	result, err := c.db.ExecContext(ctx, command, args...)
	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to execute command", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to get affected rows", err)
	}

	return affected, nil
}

// Ping checks if the database connection is still alive.
func (c *SQLConnector) Ping(ctx context.Context) error {
	if c.db == nil {
		return errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}
	return c.db.PingContext(ctx)
}

// Transaction starts a new database transaction and returns a TransactionConnector.
func (c *SQLConnector) Transaction(ctx context.Context) (TransactionConnector, error) {
	if c.db == nil {
		return nil, errors.NewError(errors.ErrorTypeDatabaseConnection, errors.ErrorMessages[errors.ErrorTypeDatabaseConnection], nil)
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeTransaction, "failed to start transaction", err)
	}

	return &SQLTransactionConnector{tx: tx}, nil
}

// SQLTransactionConnector implements the TransactionConnector interface for SQL databases.
type SQLTransactionConnector struct {
	tx *sql.Tx
}

// Query executes a query within the transaction and returns the results.
func (c *SQLTransactionConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to execute query in transaction", err)
	}
	defer rows.Close()

	
	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to get columns in transaction", err)
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to scan row in transaction", err)
		}

		row := make(map[string]interface{})
		for i, column := range columns {
			row[column] = *(values[i].(*interface{}))
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "error during row iteration in transaction", err)
	}

	return results, nil
}

// Execute executes a command within the transaction and returns the number of affected rows.
func (c *SQLTransactionConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	result, err := c.tx.ExecContext(ctx, command, args...)
	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to execute command in transaction", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.NewError(errors.ErrorTypeExecution, "failed to get affected rows in transaction", err)
	}

	return affected, nil
}

// Commit commits the transaction.
func (c *SQLTransactionConnector) Commit(ctx context.Context) error {
	if err := c.tx.Commit(); err != nil {
		return errors.NewError(errors.ErrorTypeTransaction, "failed to commit transaction", err)
	}
	return nil
}

// Rollback rolls back the transaction.
func (c *SQLTransactionConnector) Rollback(ctx context.Context) error {
	if err := c.tx.Rollback(); err != nil {
		return errors.NewError(errors.ErrorTypeTransaction, "failed to rollback transaction", err)
	}
	return nil
}