package main

import (
	"log"
	"net"

	"auth/internal/authz/service"
	"auth/internal/db"
	"auth/internal/interceptor"
	"auth/internal/repository"
	pb "auth/pb"
	"pkg/config"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the SQLite database
	client, err := db.ConnectEnt(db.Config{
		DatabaseURL: cfg.DatabaseURL,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	roleRepo := repository.NewRoleRepository(client)
	userRepo := repository.NewUserRepository(client)
	authzService := service.NewAuthzService(roleRepo, userRepo)

	lis, err := net.Listen("tcp", cfg.AuthzServiceAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	interceptor := interceptor.NewAuthInterceptor(
		interceptor.WithLogger(logger),
		// interceptor.WithTokenValidator(customTokenValidator),
		interceptor.WithSupportedSchemes(interceptor.JWT, interceptor.APIKey),
		// interceptor.WithRefreshTokenFunc(customRefreshFunc),
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor),
	)
	pb.RegisterAuthorizationServiceServer(grpcServer, authzService)

	log.Printf("Starting Authz service on %s", cfg.AuthzServiceAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
