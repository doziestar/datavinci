package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "datasource/grpc"
	"google.golang.org/grpc/reflection"

	"datasource/managers"
	"datasource/grpc/server"
)

type ServerConfig struct {
	Port int
}

func SetupAndServe(config ServerConfig) error {
	address := fmt.Sprintf(":%d", config.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	connManager := manager.NewConnectorManager()
	server := server.NewDataSourceServer(connManager)

	s := grpc.NewServer()
	pb.RegisterDataSourceServiceServer(s, server)
	
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Printf("Starting DataSource gRPC server on %s", address)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func main() {
	config := ServerConfig{
		Port: 50051,
	}

	if err := SetupAndServe(config); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}