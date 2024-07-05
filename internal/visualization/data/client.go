package client

import (
	"context"
	"encoding/json"
	"fmt"

	"datasource/managers/query"
	pb "visualization/data/client"

	"google.golang.org/grpc"
)

// DataSourceClient is a gRPC client for the DataSource service
type DataSourceClient struct {
	conn   *grpc.ClientConn
	client pb.DataSourceServiceClient
}

// NewDataSourceClient creates a new DataSourceClient
func NewDataSourceClient(address string) (*DataSourceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DataSource service: %w", err)
	}

	client := pb.NewDataSourceServiceClient(conn)
	return &DataSourceClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *DataSourceClient) Close() error {
	return c.conn.Close()
}

// Connect connects to a specific connector
func (c *DataSourceClient) Connect(ctx context.Context, connectorName string) error {
	resp, err := c.client.Connect(ctx, &pb.ConnectRequest{ConnectorName: connectorName})
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("connection failed: %s", resp.Error)
	}
	return nil
}

// Disconnect disconnects from a specific connector
func (c *DataSourceClient) Disconnect(ctx context.Context, connectorName string) error {
	resp, err := c.client.Disconnect(ctx, &pb.DisconnectRequest{ConnectorName: connectorName})
	if err != nil {
		return fmt.Errorf("failed to disconnect: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("disconnection failed: %s", resp.Error)
	}
	return nil
}

// ExecuteQuery executes a query on a specific connector
func (c *DataSourceClient) ExecuteQuery(ctx context.Context, connectorName string, q query.Query) ([]map[string]interface{}, error) {
	queryJSON, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	resp, err := c.client.ExecuteQuery(ctx, &pb.QueryRequest{
		ConnectorName: connectorName,
		Query:         string(queryJSON),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var results []map[string]interface{}
	for _, row := range resp.Rows {
		var result map[string]interface{}
		if err := json.Unmarshal(row, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal result: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// ExecuteCommand executes a command on a specific connector
func (c *DataSourceClient) ExecuteCommand(ctx context.Context, connectorName, command string, args ...string) (int64, error) {
	resp, err := c.client.ExecuteCommand(ctx, &pb.CommandRequest{
		ConnectorName: connectorName,
		Command:       command,
		Args:          args,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to execute command: %w", err)
	}
	return resp.AffectedRows, nil
}

// GetConnectors retrieves a list of available connectors
func (c *DataSourceClient) GetConnectors(ctx context.Context) ([]string, error) {
	resp, err := c.client.GetConnectors(ctx, &pb.GetConnectorsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get connectors: %w", err)
	}
	return resp.ConnectorNames, nil
}

// AddConnector adds a new connector
func (c *DataSourceClient) AddConnector(ctx context.Context, name string, config *pb.ConnectorConfig) error {
	resp, err := c.client.AddConnector(ctx, &pb.AddConnectorRequest{
		Name:   name,
		Config: config,
	})
	if err != nil {
		return fmt.Errorf("failed to add connector: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("failed to add connector: %s", resp.Error)
	}
	return nil
}

// RemoveConnector removes a connector
func (c *DataSourceClient) RemoveConnector(ctx context.Context, name string) error {
	resp, err := c.client.RemoveConnector(ctx, &pb.RemoveConnectorRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to remove connector: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("failed to remove connector: %s", resp.Error)
	}
	return nil
}
