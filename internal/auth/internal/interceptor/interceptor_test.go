package interceptor_test

import (
	"auth/internal/interceptor"
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "auth/pb"
)


func (s *mockService) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	claims, ok := interceptor.GetUserClaims(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No user claims found")
	}
	return &HelloResponse{Message: fmt.Sprintf("Hello, %v!", claims["sub"])}, nil
}

type HelloRequest struct{}
type HelloResponse struct {
	Message string
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

type mockService struct {
	pb.UnimplementedAuthServiceServer
}

func (s *mockService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	claims, ok := interceptor.GetUserClaims(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No user claims found")
	}
	return &pb.LoginResponse{AccessToken: claims["sub"].(string)}, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}


// Test helpers
func createTestServer(t *testing.T, interceptor grpc.UnaryServerInterceptor) (*grpc.Server, *bufconn.Listener) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	RegisterTestServiceServer(s, &mockService{})
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()
	return s, lis
}

func createTestClient(t *testing.T, lis *bufconn.Listener) TestServiceClient {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return NewTestServiceClient(conn)
}

// Mock JWT token generator
func generateMockToken(t *testing.T, sub string, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte("test-secret"))
	require.NoError(t, err)
	return tokenString
}

// Test cases
func TestAuthInterceptor(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	const user = "test-user"

	t.Run("Valid JWT", func(t *testing.T) {
		interceptor := interceptor.NewAuthInterceptor(
			interceptor.WithLogger(logger),
			interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": user}, nil
			}),
		)

		server, lis := createTestServer(t, interceptor)
		defer server.Stop()

		client := createTestClient(t, lis)

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer valid-token"))
		resp, err := client.SayHello(ctx, &HelloRequest{})

		assert.NoError(t, err)
		assert.Equal(t, "Hello, test-user!", resp.Message)
	})

		t.Run("Valid JWT", func(t *testing.T) {
		interceptor := interceptor.NewAuthInterceptor(
			interceptor.WithLogger(logger),
			interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": "test-user"}, nil
			}),
		)

		s := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
		pb.RegisterAuthServiceServer(s, &mockService{})
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Server exited with error: %v", err)
			}
		}()

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			t.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		md := metadata.New(map[string]string{"authorization": "Bearer valid-token"})
		ctx = metadata.NewOutgoingContext(context.Background(), md)

		resp, err := client.Login(ctx, &pb.LoginRequest{Username: "test", Password: "test"})

		require.NoError(t, err)
		assert.Equal(t, "test-user", resp.AccessToken)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		interceptor := interceptor.NewAuthInterceptor(
			interceptor.WithLogger(logger),
		)

		server, lis := createTestServer(t, interceptor)
		defer server.Stop()

		client := createTestClient(t, lis)

		_, err := client.SayHello(context.Background(), &HelloRequest{})

		assert.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("API Key Authentication", func(t *testing.T) {
		interceptor := interceptor.NewAuthInterceptor(
			interceptor.WithLogger(logger),
			interceptor.WithAPIKeyValidator(func(apiKey string) (bool, error) {
				return apiKey == "valid-api-key", nil
			}),
			interceptor.WithSupportedSchemes(interceptor.APIKey),
		)

		server, lis := createTestServer(t, interceptor)
		defer server.Stop()

		client := createTestClient(t, lis)

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "ApiKey valid-api-key"))
		resp, err := client.SayHello(ctx, &HelloRequest{})

		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "api_key")
	})

	t.Run("Token Refresh", func(t *testing.T) {
		refreshCalled := false
		interceptor := interceptor.NewAuthInterceptor(
			interceptor.WithLogger(logger),
			interceptor.WithTokenValidator(func(token string) (jwt.MapClaims, error) {
				return jwt.MapClaims{"sub": user, "exp": time.Now().Add(time.Minute).Unix()}, nil
			}),
			interceptor.WithRefreshTokenFunc(func(oldToken string) (string, error) {
				refreshCalled = true
				return "new-token", nil
			}),
			interceptor.WithTokenRefreshWindow(time.Hour), // Set a large window to ensure refresh is triggered
		)

		server, lis := createTestServer(t, interceptor)
		defer server.Stop()

		client := createTestClient(t, lis)

		ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer valid-token"))
		_, err := client.SayHello(ctx, &HelloRequest{})

		assert.NoError(t, err)
		assert.True(t, refreshCalled, "Token refresh should have been called")
	})
}

func TestTokenGenerator(t *testing.T) {
	secretKey := []byte("test-secret")
	issuer := "test-issuer"
	duration := time.Hour

	generator := interceptor.NewTokenGenerator(secretKey, issuer, duration)

	t.Run("Generate Valid Token", func(t *testing.T) {
		claims := jwt.MapClaims{"sub": "test-user"}
		token, err := generator.GenerateToken(claims)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "test-user", parsedClaims["sub"])
		assert.Equal(t, issuer, parsedClaims["iss"])
		assert.NotEmpty(t, parsedClaims["iat"])
		assert.NotEmpty(t, parsedClaims["exp"])
	})
}

func TestPasswordHasher(t *testing.T) {
	hasher := interceptor.NewPasswordHasher(10)

	t.Run("Hash and Verify Password", func(t *testing.T) {
		password := "test-password"
		hashedPassword, err := hasher.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEqual(t, password, hashedPassword)

		assert.True(t, hasher.CheckPassword(password, hashedPassword))
		assert.False(t, hasher.CheckPassword("wrong-password", hashedPassword))
	})
}

func TestBase64Encoder(t *testing.T) {
	encoder := interceptor.NewBase64Encoder()

	t.Run("Encode and Decode", func(t *testing.T) {
		original := "Hello, World!"
		encoded := encoder.Encode(original)
		decoded, err := encoder.Decode(encoded)

		assert.NoError(t, err)
		assert.NotEqual(t, original, encoded)
		assert.Equal(t, original, decoded)
	})

	t.Run("Decode Invalid Base64", func(t *testing.T) {
		_, err := encoder.Decode("invalid-base64")
		assert.Error(t, err)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("AuthMetadataKey", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer test-token"))
		key, ok := interceptor.AuthMetadataKey(ctx)
		assert.True(t, ok)
		assert.Equal(t, "Bearer test-token", key)

		_, ok = interceptor.AuthMetadataKey(context.Background())
		assert.False(t, ok)
	})

	t.Run("ExtractBearerToken", func(t *testing.T) {
		token, err := interceptor.ExtractBearerToken("Bearer test-token")
		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)

		_, err = interceptor.ExtractBearerToken("InvalidHeader test-token")
		assert.Error(t, err)
	})
}

// Mock gRPC service registration
type TestServiceServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
}

func RegisterTestServiceServer(s *grpc.Server, srv TestServiceServer) {
	s.RegisterService(&_TestService_serviceDesc, srv)
}

type TestServiceClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

type testServiceClient struct {
	cc *grpc.ClientConn
}

func NewTestServiceClient(cc *grpc.ClientConn) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/test.TestService/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _TestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "test.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _TestService_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "test_service.proto",
}

func _TestService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.TestService/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}