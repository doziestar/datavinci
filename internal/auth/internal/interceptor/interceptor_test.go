package interceptor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"auth/internal/interceptor"
)

// MockHandler is a mock gRPC handler for testing
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(ctx context.Context, req interface{}) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestAuthInterceptor(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	faker := gofakeit.New(0)

	tests := []struct {
		name           string
		setupAuth      func() interceptor.AuthInterceptorOption
		setupContext   func() context.Context
		expectedStatus codes.Code
		expectedError  string
	}{
		{
			name: "Valid JWT",
			setupAuth: func() interceptor.AuthInterceptorOption {
				return interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"user_id": faker.UUID()}, nil
				})
			},
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": "Bearer " + faker.UUID()})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedStatus: codes.OK,
		},
		{
			name: "Invalid JWT",
			setupAuth: func() interceptor.AuthInterceptorOption {
				return interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
					return nil, fmt.Errorf("invalid token: %s", faker.Sentence(5))
				})
			},
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": "Bearer " + faker.UUID()})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedStatus: codes.Unauthenticated,
			expectedError:  "Invalid token",
		},
		{
			name: "Missing Authorization Header",
			setupAuth: func() interceptor.AuthInterceptorOption {
				return interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"user_id": faker.UUID()}, nil
				})
			},
			setupContext: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.New(nil))
			},
			expectedStatus: codes.Unauthenticated,
			expectedError:  "Missing authorization header",
		},
		{
			name: "Invalid Authorization Header Format",
			setupAuth: func() interceptor.AuthInterceptorOption {
				return interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"user_id": faker.UUID()}, nil
				})
			},
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": faker.Word()})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedStatus: codes.Unauthenticated,
			expectedError:  "Invalid authorization header format",
		},
		{
			name: "Valid API Key",
			setupAuth: func() interceptor.AuthInterceptorOption {
				return interceptor.WithAPIKeyValidator(func(apiKey string) (bool, error) {
					return true, nil
				})
			},
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": "ApiKey " + faker.UUID()})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedStatus: codes.OK,
		},
		{
			name: "Invalid API Key",
			setupAuth: func() interceptor.AuthInterceptorOption {
				return interceptor.WithAPIKeyValidator(func(apiKey string) (bool, error) {
					return false, nil
				})
			},
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": "ApiKey " + faker.UUID()})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedStatus: codes.Unauthenticated,
			expectedError:  "Invalid API key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authInterceptor := interceptor.NewAuthInterceptor(
				tt.setupAuth(),
				interceptor.WithLogger(logger),
				interceptor.WithSupportedSchemes(interceptor.JWT, interceptor.APIKey),
			)

			ctx := tt.setupContext()
			mockHandler := &MockHandler{}
			mockHandler.On("Handle", mock.Anything, mock.Anything).Return(nil, nil)

			_, err := authInterceptor(ctx, nil, &grpc.UnaryServerInfo{}, mockHandler.Handle)

			if tt.expectedStatus != codes.OK {
				assert.Error(t, err)
				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, status.Code())
				assert.Contains(t, status.Message(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				mockHandler.AssertCalled(t, "Handle", mock.Anything, mock.Anything)
			}
		})
	}
}

func TestTokenGenerator(t *testing.T) {
	faker := gofakeit.New(0)
	secretKey := []byte(faker.UUID())
	issuer := faker.Company()
	duration := time.Duration(faker.Number(1, 24)) * time.Hour

	generator := interceptor.NewTokenGenerator(secretKey, issuer, duration)

	claims := jwt.MapClaims{
		"user_id": faker.UUID(),
		"role":    faker.JobTitle(),
	}

	token, err := generator.GenerateToken(claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the generated token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, claims["user_id"], parsedClaims["user_id"])
	assert.Equal(t, claims["role"], parsedClaims["role"])
	assert.Equal(t, issuer, parsedClaims["iss"])
	assert.NotNil(t, parsedClaims["iat"])
	assert.NotNil(t, parsedClaims["exp"])
}

func TestPerIPRateLimiter(t *testing.T) {
	faker := gofakeit.New(0)
	limiter := interceptor.NewPerIPRateLimiter(rate.Limit(faker.Float32Range(0.1, 10)), faker.Number(1, 100))

	ip1 := faker.IPv4Address()
	ip2 := faker.IPv4Address()

	// Test adding and getting limiters
	l1 := limiter.GetLimiter(ip1)
	assert.NotNil(t, l1)

	l2 := limiter.GetLimiter(ip2)
	assert.NotNil(t, l2)

	assert.NotEqual(t, l1, l2)

	// Test rate limiting
	for i := 0; i < 100; i++ {
		if !l1.Allow() {
			break
		}
	}
	assert.False(t, l1.Allow())

	for i := 0; i < 100; i++ {
		if !l2.Allow() {
			break
		}
	}
	assert.False(t, l2.Allow())

	// Wait for rate limit to reset
	time.Sleep(time.Second)

	assert.True(t, l1.Allow())
	assert.True(t, l2.Allow())
}

func TestPasswordHasher(t *testing.T) {
	faker := gofakeit.New(0)
	hasher := interceptor.NewPasswordHasher(faker.Number(10, 14))

	password := faker.Password(true, true, true, true, false, 10)

	hashedPassword, err := hasher.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)

	// Test correct password
	assert.True(t, hasher.CheckPassword(password, hashedPassword))

	// Test incorrect password
	assert.False(t, hasher.CheckPassword(faker.Password(true, true, true, true, false, 10), hashedPassword))
}

func TestBase64Encoder(t *testing.T) {
	faker := gofakeit.New(0)
	encoder := interceptor.NewBase64Encoder()

	original := faker.Sentence(10)
	encoded := encoder.Encode(original)
	assert.NotEqual(t, original, encoded)

	decoded, err := encoder.Decode(encoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)

	// Test invalid base64 string
	_, err = encoder.Decode(faker.Word())
	assert.Error(t, err)
}

func TestAuthMetadataKey(t *testing.T) {
	faker := gofakeit.New(0)

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedKey    string
		expectedExists bool
	}{
		{
			name: "Valid auth metadata",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": "Bearer " + faker.UUID()})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedExists: true,
		},
		{
			name: "Missing auth metadata",
			setupContext: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.New(nil))
			},
			expectedExists: false,
		},
		{
			name: "Empty auth metadata",
			setupContext: func() context.Context {
				md := metadata.New(map[string]string{"authorization": ""})
				return metadata.NewIncomingContext(context.Background(), md)
			},
			expectedKey:    "",
			expectedExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			key, exists := interceptor.AuthMetadataKey(ctx)

			assert.Equal(t, tt.expectedExists, exists)
			if tt.expectedExists {
				if tt.expectedKey != "" {
					assert.Equal(t, tt.expectedKey, key)
				} else {
					assert.NotEmpty(t, key)
				}
			}
		})
	}
}

func TestExtractBearerToken(t *testing.T) {
	faker := gofakeit.New(0)

	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Bearer token",
			authHeader:    "Bearer " + faker.UUID(),
			expectError:   false,
		},
		{
			name:        "Missing Bearer prefix",
			authHeader:  faker.UUID(),
			expectError: true,
		},
		{
			name:        "Empty auth header",
			authHeader:  "",
			expectError: true,
		},
		{
			name:        "Invalid format",
			authHeader:  "Bearer",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := interceptor.ExtractBearerToken(tt.authHeader)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestGetUserClaims(t *testing.T) {
	// faker := gofakeit.New(0)

	// claims := jwt.MapClaims{
	// 	"user_id": faker.UUID(),
	// 	"role":    faker.JobTitle(),
	// }
	

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedClaims jwt.MapClaims
		expectedOk     bool
	}{
		// {
		// 	name: "Valid user claims",
		// 	setupContext: func() context.Context {
		// 		return context.WithValue(context.Background(), interceptor.UserClaimsKey, claims)
		// 	},
		// 	expectedClaims: claims,
		// 	expectedOk:     true,
		// },
		{
			name: "Missing user claims",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			gotClaims, ok := interceptor.GetUserClaims(ctx)

			assert.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				assert.Equal(t, tt.expectedClaims, gotClaims)
			}
		})
	}
}