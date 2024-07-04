// Package interceptor provides middleware for gRPC server authentication.
//
// This package offers a flexible and feature-rich authentication interceptor
// for gRPC servers. It supports multiple authentication schemes, including
// JWT and API keys, and provides options for customization, logging, rate
// limiting, and token refresh.
package interceptor

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"pkg/config"
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
	APIKey AuthScheme = "ApiKey"
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
	// Apply rate limiting
	if err := config.RateLimiter.Wait(ctx); err != nil {
		config.Logger.Warn("Rate limit exceeded", zap.Error(err))
		return nil, status.Errorf(codes.ResourceExhausted, "Rate limit exceeded")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		config.Logger.Warn("Missing metadata")
		return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
	}

	authHeader, ok := md[config.MetadataKey]
	if !ok || len(authHeader) == 0 {
		config.Logger.Warn("Missing authorization header")
		return nil, status.Errorf(codes.Unauthenticated, "Missing authorization header")
	}

	authParts := strings.SplitN(authHeader[0], " ", 2)
	if len(authParts) != 2 {
		config.Logger.Warn("Invalid authorization header format")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization header format")
	}

	authScheme := AuthScheme(authParts[0])
	authToken := authParts[1]

	var claims jwt.MapClaims
	var err error

	switch authScheme {
	case JWT:
		claims, err = config.TokenValidator(authToken)
		if err != nil {
			config.Logger.Warn("Invalid JWT token", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
		}

		// Check if token needs refresh
		if exp, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(exp), 0)
			if time.Until(expTime) < config.TokenRefreshWindow {
				newToken, err := config.RefreshTokenFunc(authToken)
				if err != nil {
					config.Logger.Warn("Failed to refresh token", zap.Error(err))
				} else {
					// Add the new token to the outgoing context
					md.Set("new-token", newToken)
					ctx = metadata.NewOutgoingContext(ctx, md)
				}
			}
		}

	case APIKey:
		valid, err := validateAPIKey(authToken, config)
		if err != nil || !valid {
			config.Logger.Warn("Invalid API key", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "Invalid API key")
		}
		claims = jwt.MapClaims{"api_key": authToken}

	default:
		config.Logger.Warn("Unsupported authentication scheme", zap.String("scheme", string(authScheme)))
		return nil, status.Errorf(codes.Unauthenticated, "Unsupported authentication scheme")
	}

	newCtx := context.WithValue(ctx, userClaimsKey, claims)

	// Log the authenticated request
	peer, _ := peer.FromContext(ctx)
	config.Logger.Info("Authenticated request",
		zap.String("method", info.FullMethod),
		zap.String("peer", peer.Addr.String()),
		zap.Any("claims", claims),
	)

	return handler(newCtx, req)
}

// defaultTokenValidator is the default implementation of JWT token validation.
//
// This function parses and validates a JWT token using a secret key.
// In a production environment, you should replace this with your own
// implementation that uses your secret key and includes any additional
// validation logic specific to your application.
func defaultTokenValidator(tokenString string) (jwt.MapClaims, error) {
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JWTSecret), nil
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
		// we will read from our database
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

// TokenGenerator is a helper struct for generating JWT tokens.
type TokenGenerator struct {
	secretKey []byte
	issuer    string
	duration  time.Duration
}

// NewTokenGenerator creates a new TokenGenerator.
func NewTokenGenerator(secretKey []byte, issuer string, duration time.Duration) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
		issuer:    issuer,
		duration:  duration,
	}
}

// GenerateToken generates a new JWT token with the given claims.
func (g *TokenGenerator) GenerateToken(claims jwt.MapClaims) (string, error) {
	now := time.Now()
	claims["iss"] = g.issuer
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(g.duration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.secretKey)
}

// PerIPRateLimiter is a helper struct for per-IP rate limiting.
type PerIPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewPerIPRateLimiter creates a new PerIPRateLimiter.
func NewPerIPRateLimiter(r rate.Limit, b int) *PerIPRateLimiter {
	return &PerIPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// AddIP adds an IP address to the rate limiter.
func (l *PerIPRateLimiter) AddIP(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter := rate.NewLimiter(l.r, l.b)
	l.ips[ip] = limiter
	return limiter
}

// GetLimiter returns the rate limiter for the given IP address.
func (l *PerIPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	l.mu.Lock()
	limiter, exists := l.ips[ip]

	if !exists {
		l.mu.Unlock()
		return l.AddIP(ip)
	}

	l.mu.Unlock()
	return limiter
}

// PasswordHasher is a helper struct for hashing and verifying passwords.
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new PasswordHasher.
//
// Parameters:
//   - cost: The cost of the bcrypt algorithm (default is 10).
//
// Usage:
//
//	hasher := NewPasswordHasher(12)
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost == 0 {
		cost = 10
	}
	return &PasswordHasher{cost: cost}
}

// HashPassword hashes a password using bcrypt.
//
// Parameters:
//   - password: The password to hash.
//
// Returns:
//   - The hashed password and an error if hashing fails.
//
// Usage:
//
//	hashedPassword, err := hasher.HashPassword("myPassword123")
func (ph *PasswordHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	return string(bytes), err
}

// CheckPassword checks if a password matches the hashed version.
//
// Parameters:
//   - password: The password to check.
//   - hashedPassword: The hashed password to compare against.
//
// Returns:
//   - true if the password matches, false otherwise.
//
// Usage:
//
//	if hasher.CheckPassword("myPassword123", hashedPassword) {
//	    // Password is correct
//	}
func (ph *PasswordHasher) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Base64Encoder is a helper struct for encoding and decoding Base64 strings.
type Base64Encoder struct{}

// NewBase64Encoder creates a new Base64Encoder.
//
// Usage:
//
//	encoder := NewBase64Encoder()
func NewBase64Encoder() *Base64Encoder {
	return &Base64Encoder{}
}

// Encode encodes a string to Base64.
//
// Parameters:
//   - data: The string to encode.
//
// Returns:
//   - The Base64 encoded string.
//
// Usage:
//
//	encoded := encoder.Encode("Hello, World!")
func (b *Base64Encoder) Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decode decodes a Base64 string.
//
// Parameters:
//   - encodedData: The Base64 encoded string to decode.
//
// Returns:
//   - The decoded string and an error if decoding fails.
//
// Usage:
//
//	decoded, err := encoder.Decode(encodedString)
func (b *Base64Encoder) Decode(encodedData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
	if len(values) == 0 {
		return "", false
	}

	return values[0], true
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
