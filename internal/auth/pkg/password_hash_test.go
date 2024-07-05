package pkg_test

import (
	"auth/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHasherHashPassword(t *testing.T) {

	t.Run("HashPassword", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(12)
		hashedPassword, err := hasher.HashPassword("myPassword123")
		require.NoError(t, err, "HashPassword should not return an error")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	})

	t.Run("HashPasswordWithDefaultCost", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(0)
		hashedPassword, err := hasher.HashPassword("myPassword123")
		require.NoError(t, err, "HashPassword should not return an error")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	})

	t.Run("HashPasswordWithEmptyPassword", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(12)
		hashedPassword, err := hasher.HashPassword("")
		require.Error(t, err, "HashPassword should return an error")
		assert.Empty(t, hashedPassword, "Hashed password should be empty")
	})

	t.Run("HashPasswordWithInvalidCost", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(-1)
		hashedPassword, err := hasher.HashPassword("myPassword123")
		require.NoError(t, err, "HashPassword should not return an error even with invalid initial cost")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	})

	t.Run("HashPasswordWithInvalidCost", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(bcrypt.MaxCost + 1)
		hashedPassword, err := hasher.HashPassword("myPassword123")
		require.Error(t, err, "HashPassword should return an error")
		assert.Empty(t, hashedPassword, "Hashed password should be empty")
	})

	t.Run("VerifyPassword", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(12)
		password := "myPassword123"
		hashedPassword, err := hasher.HashPassword(password)
		require.NoError(t, err, "HashPassword should not return an error")

		match := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		assert.NoError(t, match, "Passwords should match")
	})

	t.Run("VerifyPasswordWithInvalidPassword", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(12)
		hashedPassword, err := hasher.HashPassword("myPassword123")
		require.NoError(t, err, "HashPassword should not return an error")

		match := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("myPassword1234"))
		assert.Error(t, match, "Passwords should not match")
	})

	t.Run("VerifyPasswordWithEmptyPassword", func(t *testing.T) {
		hasher := pkg.NewPasswordHasher(12)
		hashedPassword, err := hasher.HashPassword("myPassword123")
		require.NoError(t, err, "HashPassword should not return an error")

		match := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(""))
		assert.Error(t, match, "Passwords should not match")
	})
}
