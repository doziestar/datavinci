package pkg

import (
	"pkg/common/errors"

	"golang.org/x/crypto/bcrypt"
)

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
	if cost <= 0 {
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
	if password == "" {
		return "", errors.NewError(errors.ErrorTypeEmptyPassword, "Password cannot be empty", nil)
	}
	if ph.cost < bcrypt.MinCost || ph.cost > bcrypt.MaxCost {
		return "", errors.NewError(errors.ErrorTypeInvalidCost, "Invalid bcrypt cost", nil)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ph.cost)
	return string(bytes), err
}