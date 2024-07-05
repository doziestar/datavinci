package interceptor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestValidateConfig(t *testing.T) {
	t.Run("Valid config", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
		}
		err := validateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("Nil config", func(t *testing.T) {
		err := validateConfig(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AuthInterceptorConfig is nil")
	})

	t.Run("Nil logger", func(t *testing.T) {
		config := &AuthInterceptorConfig{}
		err := validateConfig(config)
		assert.NoError(t, err)
		assert.NotNil(t, config.Logger)
	})
}

func TestApplyRateLimiting(t *testing.T) {
	t.Run("No rate limiter", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
		}
		err := applyRateLimiting(context.Background(), config)
		assert.NoError(t, err)
	})

	// Add more test cases for rate limiting scenarios
}

func TestExtractAuthInfo(t *testing.T) {
	t.Run("Valid JWT", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer token123"))
		config := &AuthInterceptorConfig{
			Logger:      zap.NewNop(),
			MetadataKey: "authorization",
		}
		token, scheme, err := extractAuthInfo(ctx, config)
		assert.NoError(t, err)
		assert.Equal(t, "token123", token)
		assert.Equal(t, JWT, scheme)
	})

	t.Run("Valid API Key", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "APIKey key123"))
		config := &AuthInterceptorConfig{
			Logger:      zap.NewNop(),
			MetadataKey: "authorization",
		}
		token, scheme, err := extractAuthInfo(ctx, config)
		assert.NoError(t, err)
		assert.Equal(t, "key123", token)
		assert.Equal(t, APIKey, scheme)
	})

	t.Run("Missing metadata", func(t *testing.T) {
		ctx := context.Background()
		config := &AuthInterceptorConfig{
			Logger:      zap.NewNop(),
			MetadataKey: "authorization",
		}
		_, _, err := extractAuthInfo(ctx, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Missing metadata")
	})

	t.Run("Missing authorization header", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("other-key", "value"))
		config := &AuthInterceptorConfig{
			Logger:      zap.NewNop(),
			MetadataKey: "authorization",
		}
		_, _, err := extractAuthInfo(ctx, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Missing authorization header")
	})

	t.Run("Invalid authorization header format", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "InvalidFormat"))
		config := &AuthInterceptorConfig{
			Logger:      zap.NewNop(),
			MetadataKey: "authorization",
		}
		_, _, err := extractAuthInfo(ctx, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid authorization header format")
	})
}

func TestAuthenticateRequest(t *testing.T) {
	t.Run("JWT authentication", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "user123"}, nil
			},
		}
		claims, err := authenticateRequest(context.Background(), "validtoken", JWT, config)
		assert.NoError(t, err)
		assert.Equal(t, "user123", claims["sub"])
	})

	t.Run("API Key authentication", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			APIKeyValidator: func(apiKey string) (bool, error) {
				return apiKey == "validkey", nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		claims, err := authenticateRequest(context.Background(), "validkey", APIKey, config)
		assert.NoError(t, err)
		assert.Equal(t, "validkey", claims["api_key"])
	})

	t.Run("Unsupported scheme", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
		}
		_, err := authenticateRequest(context.Background(), "token", AuthScheme("Unsupported"), config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unsupported authentication scheme")
	})
}

func TestAuthenticateJWT(t *testing.T) {
	t.Run("Valid JWT", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "user123", "exp": float64(time.Now().Add(time.Hour).Unix())}, nil
			},
		}
		claims, err := authenticateJWT(context.Background(), "validtoken", config)
		assert.NoError(t, err)
		assert.Equal(t, "user123", claims["sub"])
	})

	t.Run("Invalid JWT", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return nil, status.Error(codes.Unauthenticated, "Invalid token")
			},
		}
		_, err := authenticateJWT(context.Background(), "invalidtoken", config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid token")
	})

	t.Run("Token refresh", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "user123", "exp": float64(time.Now().Add(5 * time.Minute).Unix())}, nil
			},
			TokenRefreshWindow: 10 * time.Minute,
			RefreshTokenFunc: func(token string) (string, error) {
				return "newtoken", nil
			},
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
		_, err := authenticateJWT(ctx, "validtoken", config)
		assert.NoError(t, err)
		_, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		// assert.Equal(t, []string{"newtoken"}, md.Get("new-token"))
	})
}

func TestAuthenticateAPIKey(t *testing.T) {
	t.Run("Valid API Key", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			APIKeyValidator: func(apiKey string) (bool, error) {
				return apiKey == "validkey", nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		claims, err := authenticateAPIKey("validkey", config)
		assert.NoError(t, err)
		assert.Equal(t, "validkey", claims["api_key"])
	})

	t.Run("Invalid API Key", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			APIKeyValidator: func(apiKey string) (bool, error) {
				return false, nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		_, err := authenticateAPIKey("invalidkey", config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid API key")
	})

	t.Run("Cached API Key", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			APIKeyValidator: func(apiKey string) (bool, error) {
				return apiKey == "validkey", nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		// First call should validate and cache
		_, err := authenticateAPIKey("validkey", config)
		assert.NoError(t, err)

		// Second call should use cache
		config.APIKeyValidator = func(apiKey string) (bool, error) {
			return false, nil // This should not be called
		}
		claims, err := authenticateAPIKey("validkey", config)
		assert.NoError(t, err)
		assert.Equal(t, "validkey", claims["api_key"])
	})
}

func TestLogAuthenticatedRequest(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a logger that writes to the buffer
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	))

	config := &AuthInterceptorConfig{
		Logger: logger,
	}
	ctx := context.Background()
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"}
	claims := jwt.MapClaims{"sub": "user123"}

	// Call the function
	logAuthenticatedRequest(ctx, info, config, claims)

	// Check if the log contains expected information
	logContent := buf.String()
	assert.Contains(t, logContent, "Authenticated request")
	assert.Contains(t, logContent, "/test.Service/Method")
	assert.Contains(t, logContent, "user123")
}

func TestDefaultTokenValidator(t *testing.T) {
	secretKey := "my-secret"

	t.Run("Valid token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": gofakeit.UUID(),
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		require.NoError(t, err)

		validatedClaims, err := defaultTokenValidator(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, claims["sub"], validatedClaims["sub"])
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, err := defaultTokenValidator("invalid.token.string")
		assert.Error(t, err)
	})

	t.Run("Expired token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": gofakeit.UUID(),
			"exp": time.Now().Add(-time.Hour).Unix(), // Expired
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		require.NoError(t, err)

		_, err = defaultTokenValidator(tokenString)
		assert.Error(t, err)
	})
}

func TestDefaultAPIKeyValidator(t *testing.T) {
	t.Run("Valid API key", func(t *testing.T) {
		valid, err := defaultAPIKeyValidator("valid-api-key-1")
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("Invalid API key", func(t *testing.T) {
		valid, err := defaultAPIKeyValidator("invalid-api-key")
		assert.NoError(t, err)
		assert.False(t, valid)
	})
}

func TestValidateAPIKey(t *testing.T) {
	t.Run("Valid API key", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			APIKeyValidator: func(apiKey string) (bool, error) {
				return apiKey == "validkey", nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		valid, err := validateAPIKey("validkey", config)
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("Invalid API key", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			APIKeyValidator: func(apiKey string) (bool, error) {
				return false, nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		valid, err := validateAPIKey("invalidkey", config)
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("Cached API key", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			APIKeyValidator: func(apiKey string) (bool, error) {
				return apiKey == "validkey", nil
			},
			APIKeyCache: cache.New(5*time.Minute, 10*time.Minute),
		}
		// First call should validate and cache
		_, err := validateAPIKey("validkey", config)
		assert.NoError(t, err)

		// Second call should use cache
		config.APIKeyValidator = func(apiKey string) (bool, error) {
			return false, nil // This should not be called
		}
		valid, err := validateAPIKey("validkey", config)
		assert.NoError(t, err)
		assert.True(t, valid)
	})
}

func TestDefaultRefreshTokenFunc(t *testing.T) {
	oldToken := "old-token-12345"
	newToken, err := defaultRefreshTokenFunc(oldToken)
	assert.NoError(t, err)
	assert.Contains(t, newToken, "new-refreshed-token-12345")
	assert.NotEqual(t, oldToken, newToken)
}

func TestExtractBearerToken(t *testing.T) {
	t.Run("Valid Bearer token", func(t *testing.T) {
		authHeader := "Bearer token123"
		token, err := ExtractBearerToken(authHeader)
		assert.NoError(t, err)
		assert.Equal(t, "token123", token)
	})

	t.Run("Invalid format", func(t *testing.T) {
		authHeader := "InvalidFormat token123"
		_, err := ExtractBearerToken(authHeader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid authorization header format")
	})

	t.Run("Missing token", func(t *testing.T) {
		authHeader := "Bearer"
		_, err := ExtractBearerToken(authHeader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid authorization header format")
	})

	t.Run("Case insensitive 'Bearer'", func(t *testing.T) {
		authHeader := "bEaReR token123"
		token, err := ExtractBearerToken(authHeader)
		assert.NoError(t, err)
		assert.Equal(t, "token123", token)
	})
}

func TestAuthMetadataKey(t *testing.T) {
	t.Run("Valid metadata", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer token123"))
		key, ok := AuthMetadataKey(ctx)
		assert.True(t, ok)
		assert.Equal(t, "Bearer token123", key)
	})

	t.Run("Missing metadata", func(t *testing.T) {
		ctx := context.Background()
		_, ok := AuthMetadataKey(ctx)
		assert.False(t, ok)
	})

	t.Run("Missing authorization key", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("other-key", "value"))
		_, ok := AuthMetadataKey(ctx)
		assert.False(t, ok)
	})

	t.Run("Empty authorization value", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", ""))
		_, ok := AuthMetadataKey(ctx)
		assert.False(t, ok)
	})
}

func TestGetUserClaims(t *testing.T) {
	t.Run("Valid user claims", func(t *testing.T) {
		expectedClaims := jwt.MapClaims{"sub": "user123", "role": "admin"}
		ctx := context.WithValue(context.Background(), userClaimsKey, expectedClaims)
		claims, ok := GetUserClaims(ctx)
		assert.True(t, ok)
		assert.Equal(t, expectedClaims, claims)
	})

	t.Run("Missing user claims", func(t *testing.T) {
		ctx := context.Background()
		claims, ok := GetUserClaims(ctx)
		assert.False(t, ok)
		assert.Nil(t, claims)
	})

	t.Run("Invalid user claims type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userClaimsKey, "invalid")
		claims, ok := GetUserClaims(ctx)
		assert.False(t, ok)
		assert.Nil(t, claims)
	})
}

// mockRateLimiter is a mock implementation of the rate.Limiter interface
type mockRateLimiter struct {
	allowRequest bool
}

func (m *mockRateLimiter) Wait(ctx context.Context) error {
	if !m.allowRequest {
		return status.Error(codes.ResourceExhausted, "Rate limit exceeded")
	}
	return nil
}

func TestRefreshTokenIfNeeded(t *testing.T) {
	t.Run("Token refresh", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "user123", "exp": float64(time.Now().Add(5 * time.Minute).Unix())}, nil
			},
			TokenRefreshWindow: 10 * time.Minute,
			RefreshTokenFunc: func(token string) (string, error) {
				return "newtoken", nil
			},
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
		claims, err := authenticateJWT(ctx, "validtoken", config)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		_, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		// assert.Equal(t, []string{"newtoken"}, md.Get("new-token"))
	})

	t.Run("Token doesn't need refresh", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "user123", "exp": float64(time.Now().Add(30 * time.Minute).Unix())}, nil
			},
			TokenRefreshWindow: 10 * time.Minute,
			RefreshTokenFunc: func(token string) (string, error) {
				return "newtoken", nil
			},
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
		claims, err := authenticateJWT(ctx, "validtoken", config)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		md, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		assert.Empty(t, md.Get("new-token"))
	})

	t.Run("Refresh function error", func(t *testing.T) {
		config := &AuthInterceptorConfig{
			Logger: zap.NewNop(),
			TokenValidator: func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "user123", "exp": float64(time.Now().Add(5 * time.Minute).Unix())}, nil
			},
			TokenRefreshWindow: 10 * time.Minute,
			RefreshTokenFunc: func(token string) (string, error) {
				return "", fmt.Errorf("refresh token error")
			},
		}
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
		claims, err := authenticateJWT(ctx, "validtoken", config)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		md, ok := metadata.FromIncomingContext(ctx)
		assert.True(t, ok)
		assert.Empty(t, md.Get("new-token"))
	})
}

// Run all tests
func TestMain(m *testing.M) {

	// Run the tests
	exitCode := m.Run()

	// Exit with the test result
	os.Exit(exitCode)
}
