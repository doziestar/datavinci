package pkg_test

import (
	"auth/pkg"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenGenerator(t *testing.T) {
	secretKey := []byte("test-secret-key")
	issuer := "test-issuer"
	duration := 1 * time.Hour

	t.Run("NewTokenGenerator", func(t *testing.T) {
		generator := pkg.NewTokenGenerator(secretKey, issuer, duration)
		assert.NotNil(t, generator, "NewTokenGenerator should return a non-nil generator")
	})

	t.Run("GenerateToken", func(t *testing.T) {
		generator := pkg.NewTokenGenerator(secretKey, issuer, duration)

		t.Run("BasicToken", func(t *testing.T) {
			claims := jwt.MapClaims{
				"sub":  "1234567890",
				"name": "John Doe",
			}

			token, err := generator.GenerateToken(claims)
			require.NoError(t, err, "GenerateToken should not return an error")
			assert.NotEmpty(t, token, "Generated token should not be empty")

			// Verify the token
			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			require.NoError(t, err, "Token parsing should not return an error")
			assert.True(t, parsedToken.Valid, "Token should be valid")

			parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
			require.True(t, ok, "Claims should be of type jwt.MapClaims")

			assert.Equal(t, issuer, parsedClaims["iss"], "Issuer claim should match")
			assert.Equal(t, "1234567890", parsedClaims["sub"], "Subject claim should match")
			assert.Equal(t, "John Doe", parsedClaims["name"], "Name claim should match")

			iat, ok := parsedClaims["iat"].(float64)
			require.True(t, ok, "iat claim should be a number")
			exp, ok := parsedClaims["exp"].(float64)
			require.True(t, ok, "exp claim should be a number")

			assert.InDelta(t, time.Now().Unix(), int64(iat), 5, "iat should be close to current time")
			assert.InDelta(t, time.Now().Add(duration).Unix(), int64(exp), 5, "exp should be close to current time plus duration")
		})

		t.Run("EmptyClaims", func(t *testing.T) {
			claims := jwt.MapClaims{}

			token, err := generator.GenerateToken(claims)
			require.NoError(t, err, "GenerateToken should not return an error with empty claims")
			assert.NotEmpty(t, token, "Generated token should not be empty")
		})

		t.Run("NilClaims", func(t *testing.T) {
			token, err := generator.GenerateToken(nil)
			require.NoError(t, err, "GenerateToken should not return an error with nil claims")
			assert.NotEmpty(t, token, "Generated token should not be empty")
		})

		t.Run("LargeClaims", func(t *testing.T) {
			largeClaims := jwt.MapClaims{}
			for i := 0; i < 100; i++ {
				largeClaims[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
			}

			token, err := generator.GenerateToken(largeClaims)
			require.NoError(t, err, "GenerateToken should not return an error with large claims")
			assert.NotEmpty(t, token, "Generated token should not be empty")
		})
	})
}
