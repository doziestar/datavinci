package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// validateConfig validates the configuration.
func validateConfig(config *AuthInterceptorConfig) error {
	if config == nil {
		return status.Errorf(codes.Internal, "AuthInterceptorConfig is nil")
	}
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}
	return nil
}

// applyRateLimiting applies rate limiting to the request.
func applyRateLimiting(ctx context.Context, config *AuthInterceptorConfig) error {
	if config.RateLimiter != nil {
		if err := config.RateLimiter.Wait(ctx); err != nil {
			config.Logger.Warn("Rate limit exceeded", zap.Error(err))
			return status.Errorf(codes.ResourceExhausted, "Rate limit exceeded")
		}
	}
	return nil
}

// extractAuthInfo extracts the authentication information from the context.
func extractAuthInfo(ctx context.Context, config *AuthInterceptorConfig) (string, AuthScheme, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		config.Logger.Warn("Missing metadata")
		return "", "", status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	authHeader, ok := md[config.MetadataKey]
	if !ok || len(authHeader) == 0 {
		config.Logger.Warn("Missing authorization header")
		return "", "", status.Errorf(codes.Unauthenticated, "Missing authorization header")
	}

	authParts := strings.SplitN(authHeader[0], " ", 2)
	if len(authParts) != 2 {
		config.Logger.Warn("Invalid authorization header format")
		return "", "", status.Errorf(codes.Unauthenticated, "Invalid authorization header format")
	}

	fmt.Println("authParts[1]: ", authParts[1])
	fmt.Println("authParts[0]: ", authParts[0])

	return authParts[1], AuthScheme(authParts[0]), nil
}

// authenticateRequest authenticates a request using the given token and scheme.
func authenticateRequest(ctx context.Context, authToken string, authScheme AuthScheme, config *AuthInterceptorConfig) (jwt.MapClaims, error) {
	switch authScheme {
	case JWT:
		return authenticateJWT(ctx, authToken, config)
	case APIKey:
		return authenticateAPIKey(authToken, config)
	default:
		config.Logger.Warn("Unsupported authentication scheme", zap.String("scheme", string(authScheme)))
		return nil, status.Errorf(codes.Unauthenticated, "Unsupported authentication scheme")
	}
}

// authenticateJWT authenticates a request using a JWT token.
func authenticateJWT(ctx context.Context, authToken string, config *AuthInterceptorConfig) (jwt.MapClaims, error) {
	if config.TokenValidator == nil {
		return nil, status.Errorf(codes.Internal, "TokenValidator is not configured")
	}

	claims, err := config.TokenValidator(authToken)
	if err != nil {
		config.Logger.Warn("Invalid JWT token", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	newToken, err := refreshTokenIfNeeded(ctx, authToken, claims, config)
	if err != nil {
		config.Logger.Warn("Failed to refresh token", zap.Error(err))
	} else if newToken != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}
		md = md.Copy()
		md.Set("new-token", newToken)
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	return claims, nil
}

// refreshTokenIfNeeded checks if the token needs to be refreshed and triggers a refresh if needed.
func refreshTokenIfNeeded(ctx context.Context, authToken string, claims jwt.MapClaims, config *AuthInterceptorConfig) (string, error) {
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Until(expTime) < config.TokenRefreshWindow && config.RefreshTokenFunc != nil {
			newToken, err := config.RefreshTokenFunc(authToken)
			if err != nil {
				return "", err
			}
			return newToken, nil
		}
	}
	return "", nil
}

// authenticateAPIKey authenticates a request using an API key.
func authenticateAPIKey(apiKey string, config *AuthInterceptorConfig) (jwt.MapClaims, error) {
	if config.APIKeyValidator == nil {
		return nil, status.Errorf(codes.Internal, "APIKeyValidator is not configured")
	}

	valid, err := validateAPIKey(apiKey, config)
	if err != nil || !valid {
		config.Logger.Warn("Invalid API key", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API key")
	}

	return jwt.MapClaims{"api_key": apiKey}, nil
}

// logAuthenticatedRequest logs information about an authenticated request.
//
// This function is called after a request has been successfully authenticated.
// It logs the method, peer address, and user claims.
func logAuthenticatedRequest(ctx context.Context, info *grpc.UnaryServerInfo, config *AuthInterceptorConfig, claims jwt.MapClaims) {
	if config == nil || config.Logger == nil {
		return
	}

	peerInfo := "unknown"
	if p, ok := peer.FromContext(ctx); ok && p != nil {
		peerInfo = p.Addr.String()
	}

	methodInfo := "unknown"
	if info != nil {
		methodInfo = info.FullMethod
	}

	config.Logger.Info("Authenticated request",
		zap.String("method", methodInfo),
		zap.String("peer", peerInfo),
		zap.Any("claims", claims),
	)
}

// defaultTokenValidator is the default implementation of JWT token validation.
//
// This function parses and validates a JWT token using a secret key.
// In a production environment, you should replace this with your own
// implementation that uses your secret key and includes any additional
// validation logic specific to your application.
func defaultTokenValidator(tokenString string) (jwt.MapClaims, error) {
	const secretKey = "my-secret"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// defaultAPIKeyValidator is the default implementation of API key validation.
func defaultAPIKeyValidator(apiKey string) (bool, error) {
	validKeys := map[string]bool{
		"valid-api-key-1": true,
		"key123":          true,
		"myApiKey":        true,
	}

	isValid, exists := validKeys[apiKey]
	if !exists {
		return false, nil
	}
	return isValid, nil
}

// validateAPIKey checks the validity of an API key, using caching for performance.
func validateAPIKey(apiKey string, config *AuthInterceptorConfig) (bool, error) {
	// Check cache first
	if valid, found := config.APIKeyCache.Get(apiKey); found {
		return valid.(bool), nil
	}

	// If not in cache, validate using the provided validator
	valid, err := config.APIKeyValidator(apiKey)
	if err != nil {
		return false, err
	}

	// Cache the result
	config.APIKeyCache.Set(apiKey, valid, cache.DefaultExpiration)

	return valid, nil
}

// defaultRefreshTokenFunc is the default implementation of token refresh.
func defaultRefreshTokenFunc(oldToken string) (string, error) {
	return "new-refreshed-token-" + oldToken[len(oldToken)-5:], nil
}

// ExtractBearerToken is a helper function to extract the Bearer token from an authorization header.
//
// Parameters:
//   - authHeader: The full authorization header.
//
// Returns:
//   - The Bearer token and an error if extraction fails.
//
// Usage:
//
//	token, err := ExtractBearerToken(authHeader)
//	if err != nil {
//	    return status.Errorf(codes.Unauthenticated, "Invalid authorization header")
//	}
func ExtractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return parts[1], nil
}

// AuthMetadataKey is a helper function to get the metadata key for authentication.
//
// Parameters:
//   - ctx: The context from which to extract the metadata.
//
// Returns:
//   - The authentication metadata key and a boolean indicating if it was found.
//
// Usage:
//
//	key, ok := AuthMetadataKey(ctx)
//	if !ok {
//	    return status.Errorf(codes.Unauthenticated, "Missing authentication metadata")
//	}
func AuthMetadataKey(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md.Get("authorization")
	if len(values) == 0 || values[0] == "" {
		return "", false
	}

	return values[0], true
}

// GetUserClaims retrieves the user claims from the context.
//
// This function can be used in your gRPC handlers to access the
// authenticated user's claims.
//
// Usage:
//
//	func (s *server) SomeMethod(ctx context.Context, req *pb.Request) (*pb.Response, error) {
//		claims, ok := interceptor.GetUserClaims(ctx)
//		if !ok {
//			return nil, status.Errorf(codes.Unauthenticated, "No user claims found")
//		}
//		// Use claims...
//	}
func GetUserClaims(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(jwt.MapClaims)
	return claims, ok
}
