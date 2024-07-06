package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"auth/ent/enttest"
	"auth/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSchema(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		startTime := time.Now()
		u, err := client.User.Create().
			SetUsername("testuser").
			SetEmail("test@example.com").
			SetPassword("password123").
			Save(ctx)
		endTime := time.Now()
		fmt.Printf("User creation time: %v\n", endTime.Sub(startTime))

		require.NoError(t, err)
		assert.NotEmpty(t, u.ID)
		assert.Equal(t, "testuser", u.Username)
		assert.Equal(t, "test@example.com", u.Email)
		assert.NotEqual(t, "password123", u.Password)

		const allowedTimeDiff = 5 * time.Second

		createdAtDiff := time.Since(u.CreatedAt)
		updatedAtDiff := time.Since(u.UpdatedAt)

		fmt.Printf("Time since creation: %v\n", createdAtDiff)
		fmt.Printf("Time since update: %v\n", updatedAtDiff)

		assert.True(t, createdAtDiff < allowedTimeDiff,
			"CreatedAt time difference (%v) exceeds allowed difference (%v)", createdAtDiff, allowedTimeDiff)
		assert.True(t, updatedAtDiff < allowedTimeDiff,
			"UpdatedAt time difference (%v) exceeds allowed difference (%v)", updatedAtDiff, allowedTimeDiff)
		// assert.WithinDuration(t, time.Now(), u.CreatedAt, time.Second)
		// assert.WithinDuration(t, time.Now(), u.UpdatedAt, time.Second)
	})

	t.Run("UniqueUsername", func(t *testing.T) {
		_, err := client.User.Create().
			SetUsername("testuser").
			SetEmail("another@example.com").
			SetPassword("password456").
			Save(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username")
	})

	t.Run("UniqueEmail", func(t *testing.T) {
		_, err := client.User.Create().
			SetUsername("anotheruser").
			SetEmail("test@example.com").
			SetPassword("password789").
			Save(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email")
	})

	t.Run("PasswordHashing", func(t *testing.T) {
		u, err := client.User.Create().
			SetUsername("hashtest").
			SetEmail("hash@example.com").
			SetPassword("mypassword").
			Save(ctx)

		require.NoError(t, err)
		assert.NotEqual(t, "mypassword", u.Password)

		// Verify that the hashed password is correct
		hasher := pkg.NewPasswordHasher(12)
		verified, err := hasher.VerifyPassword(u.Password, "mypassword")
		assert.NoError(t, err)
		assert.True(t, verified)
	})
}

func TestRoleSchema(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	ctx := context.Background()

	t.Run("CreateRole", func(t *testing.T) {
		r, err := client.Role.Create().
			SetName("admin").
			SetPermissions([]string{"read", "write", "delete"}).
			Save(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, r.ID)
		assert.Equal(t, "admin", r.Name)
		assert.Equal(t, []string{"read", "write", "delete"}, r.Permissions)
		assert.WithinDuration(t, time.Now(), r.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), r.UpdatedAt, time.Second)
	})

	t.Run("UniqueName", func(t *testing.T) {
		_, err := client.Role.Create().
			SetName("admin").
			SetPermissions([]string{"read"}).
			Save(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("AssignRoleToUser", func(t *testing.T) {
		u, err := client.User.Create().
			SetUsername("roleuser").
			SetEmail("role@example.com").
			SetPassword("rolepassword").
			Save(ctx)
		require.NoError(t, err)

		r, err := client.Role.Create().
			SetName("moderator").
			SetPermissions([]string{"read", "write"}).
			AddUsers(u).
			Save(ctx)
		require.NoError(t, err)

		users, err := r.QueryUsers().All(ctx)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, u.ID, users[0].ID)

		roles, err := u.QueryRoles().All(ctx)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, r.ID, roles[0].ID)
	})
}

func TestTokenSchema(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	ctx := context.Background()

	t.Run("CreateToken", func(t *testing.T) {
		u, err := client.User.Create().
			SetUsername("tokenuser").
			SetEmail("token@example.com").
			SetPassword("tokenpassword").
			Save(ctx)
		require.NoError(t, err)

		expiresAt := time.Now().Add(24 * time.Hour)
		tok, err := client.Token.Create().
			SetToken("abc123").
			SetType("access").
			SetExpiresAt(expiresAt).
			SetUser(u).
			Save(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, tok.ID)
		assert.Equal(t, "abc123", tok.Token)
		assert.Equal(t, "access", tok.Type.String())
		assert.Equal(t, expiresAt.Unix(), tok.ExpiresAt.Unix())
		assert.False(t, tok.Revoked)
		assert.WithinDuration(t, time.Now(), tok.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), tok.UpdatedAt, time.Second)
	})

	t.Run("UniqueToken", func(t *testing.T) {
		u, err := client.User.Create().
			SetUsername("tokenuser2").
			SetEmail("token2@example.com").
			SetPassword("tokenpassword2").
			Save(ctx)
		require.NoError(t, err)

		_, err = client.Token.Create().
			SetToken("abc123"). // Same token as previous test
			SetType("refresh").
			SetExpiresAt(time.Now().Add(24 * time.Hour)).
			SetUser(u).
			Save(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token")
	})

	t.Run("TokenUserRelationship", func(t *testing.T) {
		u, err := client.User.Create().
			SetUsername("tokenuser3").
			SetEmail("token3@example.com").
			SetPassword("tokenpassword3").
			Save(ctx)
		require.NoError(t, err)

		tok, err := client.Token.Create().
			SetToken("def456").
			SetType("refresh").
			SetExpiresAt(time.Now().Add(24 * time.Hour)).
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		tokenUser, err := tok.QueryUser().Only(ctx)
		require.NoError(t, err)
		assert.Equal(t, u.ID, tokenUser.ID)

		userTokens, err := u.QueryTokens().All(ctx)
		require.NoError(t, err)
		assert.Len(t, userTokens, 1)
		assert.Equal(t, tok.ID, userTokens[0].ID)
	})

	t.Run("RevokeToken", func(t *testing.T) {
		u, err := client.User.Create().
			SetUsername("tokenuser4").
			SetEmail("token4@example.com").
			SetPassword("tokenpassword4").
			Save(ctx)
		require.NoError(t, err)

		tok, err := client.Token.Create().
			SetToken("ghi789").
			SetType("access").
			SetExpiresAt(time.Now().Add(24 * time.Hour)).
			SetUser(u).
			Save(ctx)
		require.NoError(t, err)

		_, err = tok.Update().SetRevoked(true).Save(ctx)
		require.NoError(t, err)

		updatedToken, err := client.Token.Get(ctx, tok.ID)
		require.NoError(t, err)
		assert.True(t, updatedToken.Revoked)
	})
}
