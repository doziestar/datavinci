package interceptor

import (
	"context"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

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

	newCtx := context.WithValue(ctx, "user", claims)
	return handler(newCtx, req)
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	return nil, nil
}
