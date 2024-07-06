// Package connectors provides implementations for various data source connectors.
package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"pkg/common/errors"
	"pkg/common/retry"

	"github.com/gorilla/websocket"
)

// APIConnector implements the Connector interface for API data sources.
// It supports both RESTful HTTP APIs and WebSocket connections, as well as periodic polling.
type APIConnector struct {
	client            *http.Client
	config            *Config
	baseURL           string
	wsConn            *websocket.Conn
	wsBuffer          chan []byte
	pollingTicker     *time.Ticker
	stopPolling       chan struct{}
	wsLock            sync.Mutex
	reconnectBackoff  time.Duration
	maxReconnectWait  time.Duration
	stopReadWebSocket chan struct{}
}

// NewAPIConnector creates a new APIConnector with the given configuration.
//
// The config parameter should include:
//   - BaseURL: The base URL for the API
//   - TimeoutSeconds: Timeout for HTTP requests and WebSocket connection
//   - IsWebSocket: Set to true for WebSocket connections
//   - WebSocketBufferSize: Buffer size for WebSocket messages
//   - PollingIntervalSeconds: Interval for periodic polling (if > 0)
//
// Example:
//
//	config := &Config{
//	    BaseURL:                "https://api.example.com",
//	    TimeoutSeconds:         30,
//	    IsWebSocket:            false,
//	    PollingIntervalSeconds: 60,
//	}
//	connector := NewAPIConnector(config)
func NewAPIConnector(config *Config) *APIConnector {
	return &APIConnector{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.TimeoutSeconds) * time.Second,
		},
		baseURL:           config.BaseURL,
		wsBuffer:          make(chan []byte, config.WebSocketBufferSize),
		stopPolling:       make(chan struct{}),
		reconnectBackoff:  time.Second,
		maxReconnectWait:  2 * time.Minute,
		stopReadWebSocket: make(chan struct{}),
	}
}

// Connect initializes the API connector.
// For WebSocket connections, it establishes the connection.
// For polling configurations, it starts the polling routine.
//
// Example:
//
//	ctx := context.Background()
//	err := connector.Connect(ctx)
//	if err != nil {
//	    log.Fatalf("Failed to connect: %v", err)
//	}
func (c *APIConnector) Connect(ctx context.Context) error {
	if c.baseURL == "" {
		return errors.NewError(errors.ErrorTypeConfiguration, "base URL is required for API connector", nil)
	}

	if c.config.IsWebSocket {
		return c.connectWebSocket(ctx)
	}

	if c.config.PollingIntervalSeconds > 0 {
		c.startPolling(ctx)
	}

	return nil
}

// Close closes the WebSocket connection or stops the polling routine.
// It should be called when the connector is no longer needed to free up resources.
//
// Example:
//
//	ctx := context.Background()
//	err := connector.Close(ctx)
//	if err != nil {
//	    log.Printf("Error closing connector: %v", err)
//	}
func (c *APIConnector) Close(ctx context.Context) error {
	if c.wsConn != nil {
		c.stopReadWebSocket <- struct{}{}
		err := c.wsConn.Close()
		c.wsConn = nil
		return err
	}

	if c.pollingTicker != nil {
		c.stopPolling <- struct{}{}
		c.pollingTicker.Stop()
	}

	return nil
}

// Query executes a request to the API and returns the results.
// For WebSocket connections, it returns the latest message from the WebSocket.
// For HTTP connections, it sends a GET request to the specified endpoint.
//
// The query parameter is appended to the base URL for HTTP requests.
// For WebSocket connections, the query parameter is ignored.
//
// Example:
//
//	ctx := context.Background()
//	results, err := connector.Query(ctx, "/users")
//	if err != nil {
//	    log.Printf("Query failed: %v", err)
//	} else {
//	    for _, user := range results {
//	        fmt.Printf("User: %v\n", user)
//	    }
//	}
func (c *APIConnector) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if c.config.IsWebSocket {
		return c.queryWebSocket(ctx)
	}

	return c.queryHTTP(ctx, query)
}

// connectWebSocket establishes a WebSocket connection to the API.
func (c *APIConnector) connectWebSocket(ctx context.Context) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: time.Duration(c.config.TimeoutSeconds) * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, c.baseURL, nil)
	if err != nil {
		return errors.NewError(errors.ErrorTypeAPIConnection, "failed to connect to WebSocket", err)
	}

	c.wsConn = conn
	return nil
}

// readWebSocket continuously reads messages from the WebSocket connection
// and sends them to the wsBuffer channel.
func (c *APIConnector) readWebSocket() {
	defer func() {
		close(c.wsBuffer)
		close(c.stopReadWebSocket)
	}()

	for {
		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			select {
			case <-c.stopReadWebSocket:
				return
			default:
				log.Printf("WebSocket read error: %v", err)
				if err := c.reconnectWebSocket(); err != nil {
					log.Printf("Failed to reconnect WebSocket: %v", err)
					return
				}
				continue
			}
		}
		c.wsBuffer <- message
	}
}

// reconnectWebSocket attempts to reconnect the WebSocket with exponential backoff
func (c *APIConnector) reconnectWebSocket() error {
	c.wsLock.Lock()
	defer c.wsLock.Unlock()

	backoff := c.reconnectBackoff
	for {
		log.Printf("Attempting to reconnect WebSocket in %v", backoff)
		time.Sleep(backoff)

		if err := c.connectWebSocket(context.Background()); err == nil {
			log.Println("Successfully reconnected WebSocket")
			c.reconnectBackoff = time.Second // Reset backoff on successful connection
			return nil
		}

		backoff *= 2
		if backoff > c.maxReconnectWait {
			backoff = c.maxReconnectWait
		}
	}
}

// queryWebSocket returns the latest message received from the WebSocket connection.
func (c *APIConnector) queryWebSocket(ctx context.Context) ([]map[string]interface{}, error) {
	select {
	case message := <-c.wsBuffer:
		var result []map[string]interface{}
		err := json.Unmarshal(message, &result)
		if err != nil {
			return nil, errors.NewError(errors.ErrorTypeQuery, "failed to unmarshal WebSocket message", err)
		}
		return result, nil
	case <-ctx.Done():
		return nil, errors.NewError(errors.ErrorTypeQuery, "context cancelled while waiting for WebSocket message", ctx.Err())
	}
}

// queryHTTP sends an HTTP GET request to the API and returns the results.
// It uses the custom retry mechanism to handle transient errors.
func (c *APIConnector) queryHTTP(ctx context.Context, query string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, query)

	var result []map[string]interface{}
	err := retry.Retry(ctx, func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return errors.NewError(errors.ErrorTypeQuery, "failed to create request", err)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return errors.NewError(errors.ErrorTypeQuery, "failed to execute request", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.NewError(errors.ErrorTypeQuery, fmt.Sprintf("API returned non-OK status: %d", resp.StatusCode), nil)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.NewError(errors.ErrorTypeQuery, "failed to read response body", err)
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return errors.NewError(errors.ErrorTypeQuery, "failed to unmarshal response", err)
		}

		return nil
	}, retry.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return result, nil
}

// startPolling begins a periodic polling routine that calls the API at regular intervals.
func (c *APIConnector) startPolling(ctx context.Context) {
	c.pollingTicker = time.NewTicker(time.Duration(c.config.PollingIntervalSeconds) * time.Second)

	go func() {
		for {
			select {
			case <-c.pollingTicker.C:
				_, err := c.queryHTTP(ctx, "")
				if err != nil {
					log.Printf("Polling error: %v", err)
				}
			case <-c.stopPolling:
				return
			}
		}
	}()
}

// Execute sends a POST request to the API and returns the number of affected items.
// It uses the custom retry mechanism to handle transient errors.
//
// Example:
//
//	ctx := context.Background()
//	affected, err := connector.Execute(ctx, "/users/create", user)
//	if err != nil {
//	    log.Printf("Failed to create user: %v", err)
//	} else {
//	    fmt.Printf("Created %d user(s)\n", affected)
//	}
func (c *APIConnector) Execute(ctx context.Context, command string, args ...interface{}) (int64, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, command)

	var result struct {
		AffectedItems int64 `json:"affectedItems"`
	}

	err := retry.Retry(ctx, func() error {
		req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			return errors.NewError(errors.ErrorTypeExecution, "failed to create request", err)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return errors.NewError(errors.ErrorTypeExecution, "failed to execute request", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.NewError(errors.ErrorTypeExecution, fmt.Sprintf("API returned non-OK status: %d", resp.StatusCode), nil)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.NewError(errors.ErrorTypeExecution, "failed to read response body", err)
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			return errors.NewError(errors.ErrorTypeExecution, "failed to unmarshal response", err)
		}

		return nil
	}, retry.DefaultConfig())

	if err != nil {
		return 0, err
	}

	return result.AffectedItems, nil
}

// Ping checks if the API is accessible.
// For WebSocket connections, it checks if the connection is established.
// For HTTP connections, it sends a GET request to the /ping endpoint.
//
// Example:
//
//	ctx := context.Background()
//	err := connector.Ping(ctx)
//	if err != nil {
//	    log.Printf("API is not accessible: %v", err)
//	} else {
//	    fmt.Println("API is accessible")
//	}
func (c *APIConnector) Ping(ctx context.Context) error {
	if c.config.IsWebSocket {
		if c.wsConn == nil {
			return errors.NewError(errors.ErrorTypeAPIConnection, "WebSocket connection not established", nil)
		}
		return nil
	}

	_, err := c.queryHTTP(ctx, "/ping")
	return err
}

// Transaction is not supported for API connector.
// It always returns an error indicating that transactions are not supported.
func (c *APIConnector) Transaction(ctx context.Context) (TransactionConnector, error) {
	return nil, errors.NewError(errors.ErrorTypeUnsupported, "transactions are not supported for API connector", nil)
}
