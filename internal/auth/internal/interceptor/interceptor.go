// Package interceptor provides middleware for gRPC server authentication.
//
// This package offers a flexible and feature-rich authentication interceptor
// for gRPC servers. It supports multiple authentication schemes, including
// JWT and API keys, and provides options for customization, logging, rate
// limiting, and token refresh.
package interceptor

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// userClaimsKey is the context key for user claims.
const userClaimsKey contextKey = "userClaims"

// AuthScheme represents different authentication schemes.
type AuthScheme string

const (
	// JWT authentication scheme.
	JWT AuthScheme = "Bearer"
	// APIKey authentication scheme.
	APIKey AuthScheme = "APIKey"
)

// AuthInterceptorConfig holds the configuration for the AuthInterceptor.
type AuthInterceptorConfig struct {
	TokenValidator     func(string) (jwt.MapClaims, error)
	APIKeyValidator    func(string) (bool, error)
	MetadataKey        string
	Logger             *zap.Logger
	RateLimiter        *rate.Limiter
	RefreshTokenFunc   func(string) (string, error)
	SupportedSchemes   []AuthScheme
	TokenRefreshWindow time.Duration
	APIKeyCache        *cache.Cache
}

// AuthInterceptorOption is a function that modifies AuthInterceptorConfig.
type AuthInterceptorOption func(*AuthInterceptorConfig)

// WithTokenValidator sets a custom token validator function.
func WithTokenValidator(validator func(string) (jwt.MapClaims, error)) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.TokenValidator = validator
	}
}

// WithAPIKeyValidator sets a custom API key validator function.
func WithAPIKeyValidator(validator func(string) (bool, error)) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.APIKeyValidator = validator
	}
}

// WithMetadataKey sets a custom metadata key for the authorization header.
func WithMetadataKey(key string) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.MetadataKey = key
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger *zap.Logger) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.Logger = logger
	}
}

// WithRateLimiter sets a custom rate limiter.
func WithRateLimiter(limiter *rate.Limiter) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.RateLimiter = limiter
	}
}

// WithRefreshTokenFunc sets a custom refresh token function.
func WithRefreshTokenFunc(refreshFunc func(string) (string, error)) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.RefreshTokenFunc = refreshFunc
	}
}

// WithSupportedSchemes sets the supported authentication schemes.
func WithSupportedSchemes(schemes ...AuthScheme) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.SupportedSchemes = schemes
	}
}

// WithTokenRefreshWindow sets the time window before token expiration to trigger a refresh.
func WithTokenRefreshWindow(window time.Duration) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.TokenRefreshWindow = window
	}
}

// WithAPIKeyCache sets the cache for API key validation.
func WithAPIKeyCache(cache *cache.Cache) AuthInterceptorOption {
	return func(config *AuthInterceptorConfig) {
		config.APIKeyCache = cache
	}
}

// NewAuthInterceptor creates a new AuthInterceptor with the given options.
//
// Usage:
//
//	interceptor := NewAuthInterceptor(
//		WithLogger(logger),
//		WithTokenValidator(customTokenValidator),
//		WithSupportedSchemes(JWT, APIKey),
//	)
//	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
func NewAuthInterceptor(opts ...AuthInterceptorOption) grpc.UnaryServerInterceptor {
	config := &AuthInterceptorConfig{
		TokenValidator:     defaultTokenValidator,
		APIKeyValidator:    defaultAPIKeyValidator,
		MetadataKey:        "authorization",
		Logger:             zap.NewNop(),
		RateLimiter:        rate.NewLimiter(rate.Every(time.Second), 10),
		RefreshTokenFunc:   defaultRefreshTokenFunc,
		SupportedSchemes:   []AuthScheme{JWT},
		TokenRefreshWindow: 5 * time.Minute,
		APIKeyCache:        cache.New(5*time.Minute, 10*time.Minute),
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return authInterceptor(ctx, req, info, handler, config)
	}
}

// authInterceptor is the core function that performs the authentication.
func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler, config *AuthInterceptorConfig) (interface{}, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	if err := applyRateLimiting(ctx, config); err != nil {
		return nil, err
	}

	authToken, authScheme, err := extractAuthInfo(ctx, config)
	if err != nil {
		return nil, err
	}

	claims, err := authenticateRequest(ctx, authToken, authScheme, config)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, userClaimsKey, claims)
	logAuthenticatedRequest(newCtx, info, config, claims)

	return handler(newCtx, req)
}

