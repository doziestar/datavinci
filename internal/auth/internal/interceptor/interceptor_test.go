package interceptor_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"auth/internal/interceptor"
)

func TestAuthInterceptor(t *testing.T) {
	// Setup
	logger := zaptest.NewLogger(t)
	db := setupTestDatabase(t)
	defer db.Close()

	// Create a mock gRPC handler
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	tests := []struct {
		name           string
		setupAuth      func() (string, error)
		expectedCode   codes.Code
		expectedClaims jwt.MapClaims
	}{
		{
			name: "Valid JWT",
			setupAuth: func() (string, error) {
				return createTestJWT(t, "user123", time.Now().Add(1*time.Hour))
			},
			expectedCode: codes.OK,
			expectedClaims: jwt.MapClaims{
				"sub": "user123",
			},
		},
		{
			name: "Expired JWT",
			setupAuth: func() (string, error) {
				return createTestJWT(t, "user456", time.Now().Add(-1*time.Hour))
			},
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "Valid API Key",
			setupAuth: func() (string, error) {
				apiKey := gofakeit.UUID()
				err := insertAPIKey(db, apiKey)
				return apiKey, err
			},
			expectedCode: codes.OK,
			expectedClaims: jwt.MapClaims{
				"api_key": "myApiKey",
			},
		},
		{
			name: "Invalid API Key",
			setupAuth: func() (string, error) {
				return "invalid-api-key", nil
			},
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "Missing Authorization",
			setupAuth: func() (string, error) {
				return "", nil
			},
			expectedCode: codes.Unauthenticated,
		},
		{
			name: "Invalid Authorization Format",
			setupAuth: func() (string, error) {
				return "InvalidFormat", nil
			},
			expectedCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup authentication
			authToken, err := tt.setupAuth()
			require.NoError(t, err)

			// Create context with metadata
			ctx := context.Background()
			if authToken != "" {
				md := metadata.New(map[string]string{
					"authorization": fmt.Sprintf("Bearer %s", authToken),
				})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			// Create the interceptor
			newInterceptor := interceptor.NewAuthInterceptor(
				interceptor.WithLogger(logger),
				interceptor.WithTokenValidator(createTestTokenValidator(t)),
				interceptor.WithAPIKeyValidator(createTestAPIKeyValidator(db)),
				interceptor.WithSupportedSchemes(interceptor.JWT, interceptor.APIKey),
				interceptor.WithRateLimiter(rate.NewLimiter(rate.Every(time.Second), 100)),
				interceptor.WithTokenRefreshWindow(5*time.Minute),
				interceptor.WithAPIKeyCache(cache.New(5*time.Minute, 10*time.Minute)),
			)

			// Ensure the interceptor was created successfully
			require.NotNil(t, newInterceptor, "Interceptor should not be nil")

			// Call the interceptor
			_, err = newInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}, mockHandler)

			// Check the result
			if tt.expectedCode == codes.OK {
				assert.NoError(t, err)
				claims, ok := interceptor.GetUserClaims(ctx)
				assert.True(t, ok, "User claims should be present in context")
				if tt.expectedClaims["api_key"] != nil {
					assert.NotEmpty(t, claims["api_key"], "API key should not be empty")
				} else {
					assert.Equal(t, tt.expectedClaims, claims, "Claims should match expected values")
				}
			} else {
				assert.Error(t, err, "Expected an error for invalid cases")
				statusErr, ok := status.FromError(err)
				assert.True(t, ok, "Error should be a gRPC status error")
				assert.Equal(t, tt.expectedCode, statusErr.Code(), "Error code should match expected code")
			}
		})
	}
}

func createTestJWT(t *testing.T, userID string, expirationTime time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("test-secret"))
}

func createTestTokenValidator(t *testing.T) func(string) (jwt.MapClaims, error) {
	return func(tokenString string) (jwt.MapClaims, error) {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("test-secret"), nil
		})

		if err != nil {
			return nil, err
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims, nil
		}

		return nil, fmt.Errorf("invalid token")
	}
}

func createTestAPIKeyValidator(db *sql.DB) func(string) (bool, error) {
	return func(apiKey string) (bool, error) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM api_keys WHERE key = ?", apiKey).Scan(&count)
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}
}

func setupTestDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE api_keys (key TEXT PRIMARY KEY)`)
	require.NoError(t, err)

	return db
}

func insertAPIKey(db *sql.DB, apiKey string) error {
	_, err := db.Exec("INSERT INTO api_keys (key) VALUES (?)", apiKey)
	return err
}