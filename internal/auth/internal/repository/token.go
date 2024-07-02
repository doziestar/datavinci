package repository

import (
	"context"
	"time"

	"auth/ent"
	"auth/ent/token"
	"auth/ent/user"
)

type TokenRepository struct {
    client *ent.Client
}

func NewTokenRepository(client *ent.Client) *TokenRepository {
    return &TokenRepository{client: client}
}

func (r *TokenRepository) Create(ctx context.Context, userID, tokenString string, tokenType string, expiresAt time.Time) (*ent.Token, error) {
    return r.client.Token.
        Create().
        SetToken(tokenString).
        SetType(token.Type(tokenType)).
        SetExpiresAt(expiresAt).
        SetUserID(userID).
        Save(ctx)
}

func (r *TokenRepository) GetByToken(ctx context.Context, tokenString string) (*ent.Token, error) {
    return r.client.Token.
        Query().
        Where(token.Token(tokenString)).
        Only(ctx)
}

func (r *TokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
    _, err := r.client.Token.
        Update().
        Where(token.Token(tokenString)).
        SetRevoked(true).
        Save(ctx)
    return err
}

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