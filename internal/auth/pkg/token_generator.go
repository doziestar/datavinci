package pkg

import (
	"time"

	"github.com/golang-jwt/jwt"
)

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
    if claims == nil {
        claims = jwt.MapClaims{}
    }
    now := time.Now()
    claims["iss"] = g.issuer
    claims["iat"] = now.Unix()
    claims["exp"] = now.Add(g.duration).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(g.secretKey)
}