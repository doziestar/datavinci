// Package connectors provides various data source connectors for the DataVinci project.
package connectors

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"pkg/common/errors"
)

// FileConnector implements the Connector interface for file-based data sources.
// It supports reading data from JSON and CSV files.
type FileConnector struct {
	config   *Config
	basePath string
}

// NewFileConnector creates a new FileConnector with the given configuration.
//
// The config parameter should include:
//   - BasePath: The base directory path where the files are located
//
// Example:
//
//	config := &Config{
//	    BasePath: "/path/to/data/files",
//	}
//	connector := NewFileConnector(config)
func NewFileConnector(config *Config) *FileConnector {
	return &FileConnector{
		config:   config,
		basePath: config.BasePath,
	}
}

// Connect validates the base path for the file connector.
// It ensures that the specified base path exists and is accessible.
//
// Example:
//
//	ctx := context.Background()
//	err := connector.Connect(ctx)
//	if err != nil {
//	    log.Fatalf("Failed to connect: %v", err)
//	}
func (c *FileConnector) Connect(ctx context.Context) error {
	if c.basePath == "" {
		return errors.NewError(errors.ErrorTypeConfiguration, "base path is required for file connector", nil)
	}

	// Check if the base path exists and is accessible
	_, err := os.Stat(c.basePath)
	if err != nil {
		return errors.NewError(errors.ErrorTypeFileConnection, "failed to access base path", err)
	}

	return nil
}

// Close is a no-op for file connector as there's no persistent connection to close.
// It's implemented to satisfy the Connector interface.
func (c *FileConnector) Close(ctx context.Context) error {
	return nil
}

// Query reads data from a file and returns the results.
// The query parameter is treated as a relative file path from the base path.
// It supports JSON and CSV file formats.
//
// Example:
//
//	ctx := context.Background()
//	results, err := connector.Query(ctx, "users.json")
//	if err != nil {
//	    log.Printf("Query failed: %v", err)
//	} else {
//	    for _, user := range results {
//	        fmt.Printf("User: %v\n", user)
//	    }
//	}
func (c *FileConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	filePath := filepath.Join(c.basePath, query)
	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		return c.readJSONFile(filePath)
	case ".csv":
		return c.readCSVFile(filePath)
	default:
		return nil, errors.NewError(errors.ErrorTypeUnsupported, fmt.Sprintf("unsupported file type: %s", ext), nil)
	}
}

// Execute is not supported for file connector.
// It always returns an error indicating that the operation is not supported.
func (c *FileConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	return 0, errors.NewError(errors.ErrorTypeUnsupported, "execute operation is not supported for file connector", nil)
}

// Ping checks if the base path is accessible.
// It can be used to verify that the file connector is properly configured and operational.
//
// Example:
//
//	ctx := context.Background()
//	err := connector.Ping(ctx)
//	if err != nil {
//	    log.Printf("File connector is not accessible: %v", err)
//	} else {
//	    fmt.Println("File connector is accessible")
//	}
func (c *FileConnector) Ping(ctx context.Context) error {
	_, err := os.Stat(c.basePath)
	if err != nil {
		return errors.NewError(errors.ErrorTypeFileConnection, "failed to access base path", err)
	}
	return nil
}

// Transaction is not supported for file connector.
// It always returns an error indicating that transactions are not supported.
func (c *FileConnector) Transaction(ctx context.Context) (TransactionConnector, error) {
	return nil, errors.NewError(errors.ErrorTypeUnsupported, "transactions are not supported for file connector", nil)
}

// readJSONFile reads and parses a JSON file, returning the data as a slice of maps.
// It's used internally by the Query method for JSON files.
func (c *FileConnector) readJSONFile(filePath string) ([]map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to read JSON file", err)
	}

	var result []map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to unmarshal JSON data", err)
	}

	return result, nil
}

// readCSVFile reads and parses a CSV file, returning the data as a slice of maps.
// It's used internally by the Query method for CSV files.
// The first row of the CSV file is expected to contain headers.
func (c *FileConnector) readCSVFile(filePath string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to open CSV file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.NewError(errors.ErrorTypeQuery, "failed to read CSV data", err)
	}

	if len(records) == 0 {
		return nil, errors.NewError(errors.ErrorTypeQuery, "CSV file is empty", nil)
	}

	headers := records[0]
	var result []map[string]interface{}

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for i, value := range record {
			row[headers[i]] = value
		}
		result = append(result, row)
	}

	return result, nil
}
