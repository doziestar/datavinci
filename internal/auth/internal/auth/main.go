package main

import (
	"log"
	"net"
	"time"

	"auth/internal/auth/service"
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

	userRepo := repository.NewUserRepository(client)
	tokenRepo := repository.NewTokenRepository(client)
	authService := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret, time.Hour)

	lis, err := net.Listen("tcp", cfg.AuthServiceAddr)
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

	pb.RegisterAuthServiceServer(grpcServer, authService)

	log.Printf("Starting Auth service on %s", cfg.AuthServiceAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
