package repository_test

import (
	"context"
	"testing"
	"time"

	"auth/ent"
	"auth/ent/enttest"
	_ "auth/ent/token"
	"auth/internal/repository"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tokenTestSuite struct {
	client *ent.Client
	repo   repository.ITokenRepository
	faker  *gofakeit.Faker
}

func setupTokenTestSuite(t *testing.T) *tokenTestSuite {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NotNil(t, client, "Ent client should not be nil")
	
	return &tokenTestSuite{
		client: client,
		repo:   repository.NewTokenRepository(client),
		faker:  gofakeit.New(0),
	}
}

func (s *tokenTestSuite) createUser(t *testing.T) *ent.User {
	user, err := s.client.User.Create().
		SetID(s.faker.UUID()).
		SetEmail(s.faker.Email()).
		SetUsername(s.faker.Username()).
		SetPassword(s.faker.Password(true, true, true, true, false, 32)).
		Save(context.Background())
	
	require.NoError(t, err, "Failed to create user")
	return user
}

func (s *tokenTestSuite) createToken(t *testing.T, userID string, tokenString string, tokenType string, expiresAt time.Time) *ent.Token {
	token, err := s.repo.Create(context.Background(), userID, tokenString, tokenType, expiresAt)
	require.NoError(t, err, "Failed to create token")
	return token
}

func TestTokenRepository(t *testing.T) {
	suite := setupTokenTestSuite(t)
	defer suite.client.Close()

	t.Run("Create", suite.testCreate)
	t.Run("GetByToken", suite.testGetByToken)
	t.Run("RevokeToken", suite.testRevokeToken)
	t.Run("DeleteExpiredTokens", suite.testDeleteExpiredTokens)
	t.Run("GetValidTokensByUserID", suite.testGetValidTokensByUserID)
	t.Run("RevokeAllUserTokens", suite.testRevokeAllUserTokens)
}

func (s *tokenTestSuite) testCreate(t *testing.T) {
	user := s.createUser(t)
	tokenString := s.faker.UUID()
	tokenType := "access"
	expiresAt := time.Now().Add(time.Hour)

	token, err := s.repo.Create(context.Background(), user.ID, tokenString, tokenType, expiresAt)
	require.NoError(t, err, "Failed to create token")

	assert.Equal(t, tokenString, token.Token, "Token string mismatch")
	assert.Equal(t, tokenType, token.Type.String(), "Token type mismatch")
	assert.True(t, expiresAt.Equal(token.ExpiresAt), "Expiration time mismatch")

	savedToken, err := s.client.Token.Get(context.Background(), token.ID)
	require.NoError(t, err, "Failed to retrieve saved token")
	assert.Equal(t, tokenString, savedToken.Token, "Saved token string mismatch")
}

func (s *tokenTestSuite) testGetByToken(t *testing.T) {
	user := s.createUser(t)
	tokenString := s.faker.UUID()
	s.createToken(t, user.ID, tokenString, "access", time.Now().Add(time.Hour))

	retrievedToken, err := s.repo.GetByToken(context.Background(), tokenString)
	require.NoError(t, err, "Failed to get token by string")
	assert.Equal(t, tokenString, retrievedToken.Token, "Retrieved token string mismatch")

	_, err = s.repo.GetByToken(context.Background(), "non-existent-token")
	assert.Error(t, err, "Expected error when getting non-existent token")
}

func (s *tokenTestSuite) testRevokeToken(t *testing.T) {
	user := s.createUser(t)
	tokenString := s.faker.UUID()
	newToken := s.createToken(t, user.ID, tokenString, "access", time.Now().Add(time.Hour))

	err := s.repo.RevokeToken(context.Background(), newToken.Token)
	require.NoError(t, err, "Failed to revoke token")

	revokedToken, err := s.repo.GetByToken(context.Background(), newToken.Token)
	require.NoError(t, err, "Failed to get revoked token")
	assert.True(t, revokedToken.Revoked, "Token should be revoked")

	err = s.repo.RevokeToken(context.Background(), s.faker.UUID())
	assert.Error(t, err, "Expected error when revoking non-existent token")
	assert.Contains(t, err.Error(), "token not found or already revoked", "Error message should indicate token not found or already revoked")

	err = s.repo.RevokeToken(context.Background(), newToken.Token)
	assert.Error(t, err, "Expected error when revoking already revoked token")
	assert.Contains(t, err.Error(), "token not found or already revoked", "Error message should indicate token not found or already revoked")
}

func (s *tokenTestSuite) testDeleteExpiredTokens(t *testing.T) {
	user := s.createUser(t)
	expiredToken := s.createToken(t, user.ID, "expired-token", "access", time.Now().Add(-time.Hour))
	validToken := s.createToken(t, user.ID, s.faker.UUID(), "access", time.Now().Add(time.Hour))

	err := s.repo.DeleteExpiredTokens(context.Background())
	require.NoError(t, err, "Failed to delete expired tokens")

	_, err = s.client.Token.Get(context.Background(), expiredToken.ID)
	assert.Error(t, err, "Expired token should have been deleted")

	_, err = s.client.Token.Get(context.Background(), validToken.ID)
	assert.NoError(t, err, "Valid token should still exist")
}

func (s *tokenTestSuite) testGetValidTokensByUserID(t *testing.T) {
	user := s.createUser(t)
	validToken := s.createToken(t, user.ID, s.faker.UUID(), "access", time.Now().Add(time.Hour))
	s.createToken(t, user.ID, "expired-token", "access", time.Now().Add(-time.Hour))

	validTokens, err := s.repo.GetValidTokensByUserID(context.Background(), user.ID)
	require.NoError(t, err, "Failed to get valid tokens")
	assert.Len(t, validTokens, 1, "Expected 1 valid token")
	assert.Equal(t, validToken.Token, validTokens[0].Token, "Valid token string mismatch")
}

func (s *tokenTestSuite) testRevokeAllUserTokens(t *testing.T) {
	user := s.createUser(t)
	s.createToken(t, user.ID, s.faker.UUID(), "access", time.Now().Add(time.Hour))
	s.createToken(t, user.ID, s.faker.UUID(), "refresh", time.Now().Add(2*time.Hour))

	err := s.repo.RevokeAllUserTokens(context.Background(), user.ID)
	require.NoError(t, err, "Failed to revoke all user tokens")

	tokens, err := s.repo.GetValidTokensByUserID(context.Background(), user.ID)
	require.NoError(t, err, "Failed to get valid tokens")
	assert.Empty(t, tokens, "Expected 0 valid tokens after revocation")
}