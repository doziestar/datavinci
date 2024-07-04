// Package interceptor provides middleware for gRPC server authentication.
package interceptor

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"pkg/config"
)

// contextKey is a custom type for context keys.
type contextKey string

// AuthInterceptor is a gRPC server interceptor that performs JWT-based authentication.
//
// It extracts the JWT token from the "authorization" metadata, validates it,
// and adds the claims to the context for use in subsequent handlers.
//
// Usage:
//
//	server := grpc.NewServer(
//		grpc.UnaryInterceptor(interceptor.AuthInterceptor),
//	)
//
// The handler can then access the user claims like this:
//
//	func (s *server) SomeMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
//		claims, ok := ctx.Value("user").(jwt.MapClaims)
//		if !ok {
//			return nil, status.Errorf(codes.Unauthenticated, "No user claims found")
//		}
//		
//	}
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Missing authorization header")
	}

	token := authHeader[0]
	claims, err := validateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}
	
	const userKey contextKey = "user"
	
	newCtx := context.WithValue(ctx, userKey, claims)
	return handler(newCtx, req)
}

// validateToken verifies the given JWT token and returns its claims.
//
// It uses the JWT secret from the configuration to validate the token.
//
// Returns:
//   - jwt.MapClaims: The claims of the token if it is valid.
//   - error: An error if the token is invalid.
func validateToken(tokenString string) (jwt.MapClaims, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
