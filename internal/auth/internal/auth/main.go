package main

import (
    "log"
    "net"
    "time"

    "pkg/config"
    "auth/internal/db"
    "auth/internal/repository"
    "auth/internal/auth/service"
    pb "auth/pb"

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
    authService := service.NewAuthService(*userRepo, tokenRepo, cfg.JWTSecret, time.Hour)

    lis, err := net.Listen("tcp", cfg.AuthServiceAddr)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterAuthServiceServer(grpcServer, authService)

    log.Printf("Starting Auth service on %s", cfg.AuthServiceAddr)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}