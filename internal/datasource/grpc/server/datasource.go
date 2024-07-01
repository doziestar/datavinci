package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"datasource/grpc"
	"datasource/managers"
	"datasource/managers/query"
	"datasource/connectors"
)

type DataSourceServer struct {
	grpc.UnimplementedDataSourceServiceServer
	manager *manager.ConnectorManager
}

func NewDataSourceServer(manager *manager.ConnectorManager) *DataSourceServer {
	return &DataSourceServer{manager: manager}
}

func (s *DataSourceServer) Connect(ctx context.Context, req *grpc.ConnectRequest) (*grpc.ConnectResponse, error) {
	log.Printf("Received Connect request for connector: %s", req.ConnectorName)
	
	connector, err := s.manager.GetConnector(req.ConnectorName)
	if err != nil {
		log.Printf("Error getting connector %s: %v", req.ConnectorName, err)
		return &grpc.ConnectResponse{Success: false, Error: err.Error()}, nil
	}

	err = connector.Connect(ctx)
	if err != nil {
		log.Printf("Error connecting to %s: %v", req.ConnectorName, err)
		return &grpc.ConnectResponse{Success: false, Error: err.Error()}, nil
	}

	log.Printf("Successfully connected to %s", req.ConnectorName)
	return &grpc.ConnectResponse{Success: true}, nil
}

func (s *DataSourceServer) Disconnect(ctx context.Context, req *grpc.DisconnectRequest) (*grpc.DisconnectResponse, error) {
	log.Printf("Received Disconnect request for connector: %s", req.ConnectorName)
	
	connector, err := s.manager.GetConnector(req.ConnectorName)
	if err != nil {
		log.Printf("Error getting connector %s: %v", req.ConnectorName, err)
		return &grpc.DisconnectResponse{Success: false, Error: err.Error()}, nil
	}

	err = connector.Close(ctx)
	if err != nil {
		log.Printf("Error disconnecting from %s: %v", req.ConnectorName, err)
		return &grpc.DisconnectResponse{Success: false, Error: err.Error()}, nil
	}

	log.Printf("Successfully disconnected from %s", req.ConnectorName)
	return &grpc.DisconnectResponse{Success: true}, nil
}

func (s *DataSourceServer) ExecuteQuery(ctx context.Context, req *grpc.QueryRequest) (*grpc.QueryResponse, error) {
	log.Printf("Received ExecuteQuery request for connector: %s", req.ConnectorName)
	
	connector, err := s.manager.GetConnector(req.ConnectorName)
	if err != nil {
		log.Printf("Error getting connector %s: %v", req.ConnectorName, err)
		return nil, status.Errorf(codes.NotFound, "connector not found: %v", err)
	}

	executor := query.NewQueryExecutor(connector)
	
	var q query.Query
	err = json.Unmarshal([]byte(req.Query), &q)
	if err != nil {
		log.Printf("Error unmarshalling query: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid query: %v", err)
	}

	results, err := executor.Execute(ctx, q)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, status.Errorf(codes.Internal, "query execution failed: %v", err)
	}

	var rows [][]byte
	for _, result := range results {
		rowBytes, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error marshalling result: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to marshal result: %v", err)
		}
		rows = append(rows, rowBytes)
	}

	log.Printf("Successfully executed query on %s, returned %d rows", req.ConnectorName, len(rows))
	return &grpc.QueryResponse{Rows: rows}, nil
}

func (s *DataSourceServer) ExecuteCommand(ctx context.Context, req *grpc.CommandRequest) (*grpc.CommandResponse, error) {
	log.Printf("Received ExecuteCommand request for connector: %s", req.ConnectorName)
	
	connector, err := s.manager.GetConnector(req.ConnectorName)
	if err != nil {
		log.Printf("Error getting connector %s: %v", req.ConnectorName, err)
		return nil, status.Errorf(codes.NotFound, "connector not found: %v", err)
	}

	affected, err := connector.Execute(ctx, req.Command, req.Args)
	if err != nil {
		log.Printf("Error executing command on %s: %v", req.ConnectorName, err)
		return nil, status.Errorf(codes.Internal, "command execution failed: %v", err)
	}

	log.Printf("Successfully executed command on %s, affected %d rows", req.ConnectorName, affected)
	return &grpc.CommandResponse{AffectedRows: affected}, nil
}

// func (s *DataSourceServer) GetConnectors(ctx context.Context, req *grpc.GetConnectorsRequest) (*grpc.GetConnectorsResponse, error) {
// 	log.Printf("Received GetConnectors request")
	
// 	connectorNames, err := s.manager.GetConnector()
// 	log.Printf("Retrieved %d connectors", len(connectorNames))
// 	return &grpc.GetConnectorsResponse{ConnectorNames: connectorNames}, nil
// }

func (s *DataSourceServer) AddConnector(ctx context.Context, req *grpc.AddConnectorRequest) (*grpc.AddConnectorResponse, error) {
	log.Printf("Received AddConnector request for connector: %s", req.Name)
	
	config := &connectors.Config{
		Type:     req.Config.Type,
		Host:     req.Config.Host,
		Port:     int(req.Config.Port),
		Username: req.Config.Username,
		Password: req.Config.Password,
		Database: req.Config.Database,
	}

	for k, v := range req.Config.Options {
		if config.Options == nil {
			config.Options = make(map[string]interface{})
		}
		config.Options[k] = v
	}

	err := s.manager.AddConnector(req.Name, config)
	if err != nil {
		log.Printf("Error adding connector %s: %v", req.Name, err)
		return &grpc.AddConnectorResponse{Success: false, Error: err.Error()}, nil
	}

	log.Printf("Successfully added connector: %s", req.Name)
	return &grpc.AddConnectorResponse{Success: true}, nil
}

func (s *DataSourceServer) RemoveConnector(ctx context.Context, req *grpc.RemoveConnectorRequest) (*grpc.RemoveConnectorResponse, error) {
	log.Printf("Received RemoveConnector request for connector: %s", req.Name)
	
	err := s.manager.RemoveConnector(req.Name)
	if err != nil {
		log.Printf("Error removing connector %s: %v", req.Name, err)
		return &grpc.RemoveConnectorResponse{Success: false, Error: err.Error()}, nil
	}

	log.Printf("Successfully removed connector: %s", req.Name)
	return &grpc.RemoveConnectorResponse{Success: true}, nil
}

func (s *DataSourceServer) validateConnectorConfig(config *grpc.ConnectorConfig) error {
	if config.Type == "" {
		return fmt.Errorf("connector type is required")
	}
	if config.Host == "" {
		return fmt.Errorf("host is required")
	}
	if config.Port == 0 {
		return fmt.Errorf("port is required")
	}
	// Add more validation as needed
	return nil
}