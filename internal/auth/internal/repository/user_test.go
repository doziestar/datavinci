package repository_test

import (
	"context"
	"fmt"
	"testing"

	"auth/ent"
	"auth/ent/enttest"
	"auth/internal/repository"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type userTestSuite struct {
	client *ent.Client
	repo   repository.IUserRepository
	faker  *gofakeit.Faker
}

func setupUserTestSuite(t *testing.T) *userTestSuite {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NotNil(t, client, "Ent client should not be nil")

	return &userTestSuite{
		client: client,
		repo:   repository.NewUserRepository(client),
		faker:  gofakeit.New(0),
	}
}

func (s *userTestSuite) createUser(t *testing.T) *ent.User {
	user, err := s.repo.Create(context.Background(), &ent.User{
		Username: s.faker.Username(),
		Email:    s.faker.Email(),
		Password: s.faker.Password(true, true, true, true, false, 32),
	})
	require.NoError(t, err, "Failed to create user")
	return user
}

func (s *userTestSuite) createUniqueRole(t *testing.T) *ent.Role {
	ctx := context.Background()
	roleName := fmt.Sprintf("role_%s", s.faker.UUID())
	role, err := s.client.Role.Create().
		SetName(roleName).
		SetPermissions([]string{"read", "write"}).
		Save(ctx)
	require.NoError(t, err, "Failed to create role")
	return role
}

func TestUserRepository(t *testing.T) {
	suite := setupUserTestSuite(t)
	defer suite.client.Close()

	t.Run("Create", suite.testCreate)
	t.Run("GetByID", suite.testGetByID)
	t.Run("GetByUsername", suite.testGetByUsername)
	t.Run("GetByEmail", suite.testGetByEmail)
	t.Run("Update", suite.testUpdate)
	t.Run("Delete", suite.testDelete)
	t.Run("List", suite.testList)
	t.Run("Count", suite.testCount)
	t.Run("AddRole", suite.testAddRole)
	t.Run("RemoveRole", suite.testRemoveRole)
	t.Run("GetRoles", suite.testGetRoles)
	t.Run("Search", suite.testSearch)
	t.Run("ChangePassword", suite.testChangePassword)
	t.Run("GetUsersByRole", suite.testGetUsersByRole)
	t.Run("CheckPassword", suite.testCheckPassword)
	t.Run("SetPassword", suite.testSetPassword)
}

func (s *userTestSuite) testCreate(t *testing.T) {
	ctx := context.Background()
	newUser := &ent.User{
		Username: s.faker.Username(),
		Email:    s.faker.Email(),
		Password: s.faker.Password(true, true, true, true, false, 32),
	}

	createdUser, err := s.repo.Create(ctx, newUser)
	require.NoError(t, err, "Failed to create user")
	assert.Equal(t, newUser.Username, createdUser.Username, "Username mismatch")
	assert.Equal(t, newUser.Email, createdUser.Email, "Email mismatch")

	savedUser, err := s.client.User.Get(ctx, createdUser.ID)
	require.NoError(t, err, "Failed to retrieve saved user")
	assert.Equal(t, newUser.Username, savedUser.Username, "Saved username mismatch")
}

func (s *userTestSuite) testGetByID(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)

	retrievedUser, err := s.repo.GetByID(ctx, user.ID)
	require.NoError(t, err, "Failed to get user by ID")
	assert.Equal(t, user.ID, retrievedUser.ID, "Retrieved user ID mismatch")

	_, err = s.repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err, "Expected error when getting non-existent user")
}

func (s *userTestSuite) testGetByUsername(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)

	retrievedUser, err := s.repo.GetByUsername(ctx, user.Username)
	require.NoError(t, err, "Failed to get user by username")
	assert.Equal(t, user.Username, retrievedUser.Username, "Retrieved username mismatch")

	_, err = s.repo.GetByUsername(ctx, "non-existent-username")
	assert.Error(t, err, "Expected error when getting non-existent username")
}

func (s *userTestSuite) testGetByEmail(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)

	retrievedUser, err := s.repo.GetByEmail(ctx, user.Email)
	require.NoError(t, err, "Failed to get user by email")
	assert.Equal(t, user.Email, retrievedUser.Email, "Retrieved email mismatch")

	_, err = s.repo.GetByEmail(ctx, "non-existent@example.com")
	assert.Error(t, err, "Expected error when getting non-existent email")
}

func (s *userTestSuite) testUpdate(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)

	updatedUser := user
	updatedUser.Email = s.faker.Email()

	result, err := s.repo.Update(ctx, updatedUser)
	require.NoError(t, err, "Failed to update user")
	assert.Equal(t, updatedUser.Email, result.Email, "Updated email mismatch")

	retrievedUser, err := s.repo.GetByID(ctx, user.ID)
	require.NoError(t, err, "Failed to retrieve updated user")
	assert.Equal(t, updatedUser.Email, retrievedUser.Email, "Retrieved updated email mismatch")
}

func (s *userTestSuite) testDelete(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)

	err := s.repo.Delete(ctx, user.ID)
	require.NoError(t, err, "Failed to delete user")

	_, err = s.repo.GetByID(ctx, user.ID)
	assert.Error(t, err, "Expected error when getting deleted user")

	err = s.repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err, "Expected error when deleting non-existent user")
}

func (s *userTestSuite) testList(t *testing.T) {
	ctx := context.Background()

	_, err := s.client.User.Delete().Exec(ctx)
	require.NoError(t, err, "Failed to clear existing users")

	for i := 0; i < 5; i++ {
		s.createUser(t)
	}

	users, err := s.repo.List(ctx, 0, 3)
	require.NoError(t, err, "Failed to list users")
	assert.Len(t, users, 3, "Expected 3 users")

	users, err = s.repo.List(ctx, 3, 3)
	require.NoError(t, err, "Failed to list users")
	assert.Len(t, users, 2, "Expected 2 users")
}

func (s *userTestSuite) testCount(t *testing.T) {
	ctx := context.Background()

	initialCount, err := s.repo.Count(ctx)
	require.NoError(t, err, "Failed to count users")

	for i := 0; i < 5; i++ {
		s.createUser(t)
	}

	count, err := s.repo.Count(ctx)
	require.NoError(t, err, "Failed to count users")
	assert.Equal(t, initialCount+5, count, "Expected count to increase by 5")
}

func (s *userTestSuite) testAddRole(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)
	role := s.createUniqueRole(t)

	err := s.repo.AddRole(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to add role to user")

	roles, err := s.repo.GetRoles(ctx, user.ID)
	require.NoError(t, err, "Failed to get user roles")
	assert.Len(t, roles, 1, "Expected 1 role")
	assert.Equal(t, role.ID, roles[0].ID, "Role ID mismatch")
}

func (s *userTestSuite) testRemoveRole(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)
	role := s.createUniqueRole(t)

	err := s.repo.AddRole(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to add role to user")

	err = s.repo.RemoveRole(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to remove role from user")

	roles, err := s.repo.GetRoles(ctx, user.ID)
	require.NoError(t, err, "Failed to get user roles")
	assert.Empty(t, roles, "Expected no roles")
}

func (s *userTestSuite) testGetRoles(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)
	role := s.createUniqueRole(t)

	err := s.repo.AddRole(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to add role to user")

	roles, err := s.repo.GetRoles(ctx, user.ID)
	require.NoError(t, err, "Failed to get user roles")
	assert.Len(t, roles, 1, "Expected 1 role")
	assert.Equal(t, role.ID, roles[0].ID, "Role ID mismatch")
}

func (s *userTestSuite) testGetUsersByRole(t *testing.T) {
	ctx := context.Background()
	role := s.createUniqueRole(t)

	user1 := s.createUser(t)
	user2 := s.createUser(t)

	err := s.repo.AddRole(ctx, user1.ID, role.ID)
	require.NoError(t, err, "Failed to add role to user1")
	err = s.repo.AddRole(ctx, user2.ID, role.ID)
	require.NoError(t, err, "Failed to add role to user2")

	usersWithRole, err := s.repo.GetUsersByRole(ctx, role.ID)
	require.NoError(t, err, "Failed to get users by role")
	assert.Len(t, usersWithRole, 2, "Expected 2 users with the role")
}

func (s *userTestSuite) testCheckPassword(t *testing.T) {
	ctx := context.Background()
	password := s.faker.Password(true, true, true, true, false, 32)
	username := s.faker.Username()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err, "Failed to hash password")

	user, err := s.repo.Create(ctx, &ent.User{
		Username: username,
		Email:    s.faker.Email(),
		Password: string(hashedPassword),
	})
	require.NoError(t, err, "Failed to create user")

	isCorrect := s.repo.CheckPassword(ctx, username, password)
	assert.True(t, isCorrect, "Password check should succeed")

	isCorrect = s.repo.CheckPassword(ctx, username, "wrongPassword")
	assert.False(t, isCorrect, "Password check should fail for incorrect password")

	isCorrect = s.repo.CheckPassword(ctx, username, user.Password)
	assert.False(t, isCorrect, "Password check should fail for hashed password")
}

func (s *userTestSuite) testSetPassword(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)
	oldPassword := user.Password
	newPassword := s.faker.Password(true, true, true, true, false, 32)

	err := s.repo.SetPassword(ctx, user.Username, newPassword)
	require.NoError(t, err, "Failed to set new password")

	isCorrect := s.repo.CheckPassword(ctx, user.Username, newPassword)
	assert.True(t, isCorrect, "New password check should succeed")

	isCorrect = s.repo.CheckPassword(ctx, user.Username, oldPassword)
	assert.False(t, isCorrect, "Old password check should fail")
}

func (s *userTestSuite) testSearch(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)

	results, err := s.repo.Search(ctx, user.Username[:3])
	require.NoError(t, err, "Failed to search users")
	assert.NotEmpty(t, results, "Expected search results")

	found := false
	for _, result := range results {
		if result.ID == user.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find the created user in search results")
}

func (s *userTestSuite) testChangePassword(t *testing.T) {
	ctx := context.Background()
	user := s.createUser(t)
	newPassword := "newSecurePassword123"

	err := s.repo.ChangePassword(ctx, user.ID, newPassword)
	require.NoError(t, err, "Failed to change password")

	updatedUser, err := s.repo.GetByID(ctx, user.ID)
	require.NoError(t, err, "Failed to get updated user")
	assert.NotEqual(t, user.Password, updatedUser.Password, "Password should have changed")
}
