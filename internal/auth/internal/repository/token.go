// Package repository provides interfaces and implementations for token management.
package repository

import (
	"auth/ent"
	"auth/ent/token"
	"auth/ent/user"
	"context"
	"time"
)

// ITokenRepository defines the interface for token-related operations.
type ITokenRepository interface {
	// Create creates a new token.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - userID: The ID of the user associated with the token.
	//   - tokenString: The unique string representation of the token.
	//   - tokenType: The type of the token (e.g., "access", "refresh").
	//   - expiresAt: The expiration time of the token.
	//
	// Returns:
	//   - *ent.Token: The created token entity.
	//   - error: An error if the creation fails, nil otherwise.
	Create(ctx context.Context, userID, tokenString string, tokenType string, expiresAt time.Time) (*ent.Token, error)

	// GetByToken retrieves a token by its string representation.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - tokenString: The unique string representation of the token to retrieve.
	//
	// Returns:
	//   - *ent.Token: The retrieved token entity.
	//   - error: An error if the token is not found or if the retrieval fails, nil otherwise.
	GetByToken(ctx context.Context, tokenString string) (*ent.Token, error)

	// RevokeToken marks a token as revoked.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - tokenString: The unique string representation of the token to revoke.
	//
	// Returns:
	//   - error: An error if the revocation fails, nil otherwise.
	RevokeToken(ctx context.Context, tokenString string) error

	// DeleteExpiredTokens removes all expired or revoked tokens from the database.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//
	// Returns:
	//   - error: An error if the deletion fails, nil otherwise.
	DeleteExpiredTokens(ctx context.Context) error

	// GetValidTokensByUserID retrieves all valid tokens for a specific user.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - userID: The ID of the user whose tokens to retrieve.
	//
	// Returns:
	//   - []*ent.Token: A slice of valid token entities.
	//   - error: An error if the retrieval fails, nil otherwise.
	GetValidTokensByUserID(ctx context.Context, userID string) ([]*ent.Token, error)

	// RevokeAllUserTokens revokes all tokens belonging to a specific user.
	//
	// Parameters:
	//   - ctx: The context for the database operation.
	//   - userID: The ID of the user whose tokens to revoke.
	//
	// Returns:
	//   - error: An error if the revocation fails, nil otherwise.
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

// TokenRepository implements the ITokenRepository interface.
type TokenRepository struct {
	client *ent.Client
}

// NewTokenRepository creates a new instance of TokenRepository.
//
// Parameters:
//   - client: The Ent ORM client for database operations.
//
// Returns:
//   - *TokenRepository: A new instance of TokenRepository.
func NewTokenRepository(client *ent.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

// Create implements ITokenRepository.Create.
func (r *TokenRepository) Create(ctx context.Context, userID, tokenString string, tokenType string, expiresAt time.Time) (*ent.Token, error) {
	return r.client.Token.
		Create().
		SetToken(tokenString).
		SetType(token.Type(tokenType)).
		SetExpiresAt(expiresAt).
		SetUserID(userID).
		Save(ctx)
}

// GetByToken implements ITokenRepository.GetByToken.
func (r *TokenRepository) GetByToken(ctx context.Context, tokenString string) (*ent.Token, error) {
	return r.client.Token.
		Query().
		Where(token.Token(tokenString)).
		Only(ctx)
}

// RevokeToken implements ITokenRepository.RevokeToken.
func (r *TokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
	_, err := r.client.Token.
		Update().
		Where(token.Token(tokenString)).
		SetRevoked(true).
		Save(ctx)
	return err
}

// DeleteExpiredTokens implements ITokenRepository.DeleteExpiredTokens.
func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	_, err := r.client.Token.
		Delete().
		Where(
			token.Or(
				token.ExpiresAtLT(time.Now()),
				token.Revoked(true),
			),
		).
		Exec(ctx)
	return err
}

// GetValidTokensByUserID implements ITokenRepository.GetValidTokensByUserID.
func (r *TokenRepository) GetValidTokensByUserID(ctx context.Context, userID string) ([]*ent.Token, error) {
	return r.client.Token.
		Query().
		Where(
			token.HasUserWith(user.ID(userID)),
			token.ExpiresAtGT(time.Now()),
			token.Revoked(false),
		).
		All(ctx)
}

// RevokeAllUserTokens implements ITokenRepository.RevokeAllUserTokens.
func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	_, err := r.client.Token.
		Update().
		Where(
			token.HasUserWith(user.ID(userID)),
			token.Revoked(false),
		).
		SetRevoked(true).
		Save(ctx)
	return err
}
