package interceptor

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestNewAuthInterceptor(t *testing.T) {
	interceptor := NewAuthInterceptor()
	assert.NotNil(t, interceptor)
}

func TestWithOptions(t *testing.T) {
	logger := zap.NewExample()
	tokenValidator := func(string) (jwt.MapClaims, error) { return nil, nil }
	apiKeyValidator := func(string) (bool, error) { return true, nil }
	refreshFunc := func(string) (string, error) { return "", nil }
	cache := cache.New(5*time.Minute, 10*time.Minute)

	interceptor := NewAuthInterceptor(
		WithLogger(logger),
		WithTokenValidator(tokenValidator),
		WithAPIKeyValidator(apiKeyValidator),
		WithMetadataKey("custom-auth"),
		WithRateLimiter(rate.NewLimiter(rate.Every(time.Second), 5)),
		WithRefreshTokenFunc(refreshFunc),
		WithSupportedSchemes(JWT, APIKey),
		WithTokenRefreshWindow(10*time.Minute),
		WithAPIKeyCache(cache),
	)

	assert.NotNil(t, interceptor)
}

func TestAuthInterceptor(t *testing.T) {
	config := &AuthInterceptorConfig{
		TokenValidator: func(token string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": "user123"}, nil
		},
		APIKeyValidator: func(apiKey string) (bool, error) {
			return apiKey == "valid-api-key", nil
		},
		MetadataKey:      "authorization",
		Logger:           zap.NewNop(),
		SupportedSchemes: []AuthScheme{JWT, APIKey},
		APIKeyCache:      cache.New(5*time.Minute, 10*time.Minute), // Add this line
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	t.Run("Valid JWT", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer valid-token"))
		resp, err := authInterceptor(ctx, "request", info, handler, config)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("Valid API Key", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "APIKey valid-api-key"))
		resp, err := authInterceptor(ctx, "request", info, handler, config)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("Invalid Auth", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Invalid auth"))
		_, err := authInterceptor(ctx, "request", info, handler, config)
		assert.Error(t, err)
	})

	t.Run("Missing Auth", func(t *testing.T) {
		ctx := context.Background()
		_, err := authInterceptor(ctx, "request", info, handler, config)
		assert.Error(t, err)
	})
}
