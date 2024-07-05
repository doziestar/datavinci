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
)

type roleTestSuite struct {
	client *ent.Client
	repo   repository.IRoleRepository
	faker  *gofakeit.Faker
}

func setupRoleTestSuite(t *testing.T) *roleTestSuite {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	require.NotNil(t, client, "Ent client should not be nil")

	return &roleTestSuite{
		client: client,
		repo:   repository.NewRoleRepository(client),
		faker:  gofakeit.New(0),
	}
}

func (s *roleTestSuite) createRole(t *testing.T) *ent.Role {
	role, err := s.repo.Create(context.Background(), &ent.Role{
		Name:        fmt.Sprintf("role_%s", s.faker.UUID()),
		Permissions: []string{"read", "write"},
	})
	require.NoError(t, err, "Failed to create role")
	return role
}

func TestRoleRepository(t *testing.T) {
	suite := setupRoleTestSuite(t)
	defer suite.client.Close()

	t.Run("Create", suite.testCreate)
	t.Run("GetByID", suite.testGetByID)
	t.Run("GetByName", suite.testGetByName)
	t.Run("Update", suite.testUpdate)
	t.Run("Delete", suite.testDelete)
	t.Run("List", suite.testList)
	t.Run("Count", suite.testCount)
	t.Run("AddPermission", suite.testAddPermission)
	t.Run("RemovePermission", suite.testRemovePermission)
	t.Run("GetUsersInRole", suite.testGetUsersInRole)
	t.Run("Search", suite.testSearch)
	t.Run("GetRolesByUserID", suite.testGetRolesByUserID)
	t.Run("AssignRoleToUser", suite.testAssignRoleToUser)
	t.Run("RemoveRoleFromUser", suite.testRemoveRoleFromUser)
}

func (s *roleTestSuite) testCreate(t *testing.T) {
	ctx := context.Background()
	newRole := &ent.Role{
		Name:        fmt.Sprintf("role_%s", s.faker.UUID()),
		Permissions: []string{"read", "write"},
	}

	createdRole, err := s.repo.Create(ctx, newRole)
	require.NoError(t, err, "Failed to create role")
	assert.Equal(t, newRole.Name, createdRole.Name, "Role name mismatch")
	assert.ElementsMatch(t, newRole.Permissions, createdRole.Permissions, "Role permissions mismatch")
}

func (s *roleTestSuite) testGetByID(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)

	fetchedRole, err := s.repo.GetByID(ctx, role.ID)
	require.NoError(t, err, "Failed to get role by ID")
	assert.Equal(t, role.ID, fetchedRole.ID, "Role ID mismatch")
	assert.Equal(t, role.Name, fetchedRole.Name, "Role name mismatch")
	assert.ElementsMatch(t, role.Permissions, fetchedRole.Permissions, "Role permissions mismatch")
}

func (s *roleTestSuite) testGetByName(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)

	fetchedRole, err := s.repo.GetByName(ctx, role.Name)
	require.NoError(t, err, "Failed to get role by name")
	assert.Equal(t, role.ID, fetchedRole.ID, "Role ID mismatch")
	assert.Equal(t, role.Name, fetchedRole.Name, "Role name mismatch")
	assert.ElementsMatch(t, role.Permissions, fetchedRole.Permissions, "Role permissions mismatch")
}

func (s *roleTestSuite) testUpdate(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)

	role.Permissions = append(role.Permissions, "delete")
	updatedRole, err := s.repo.Update(ctx, role)
	require.NoError(t, err, "Failed to update role")
	assert.ElementsMatch(t, role.Permissions, updatedRole.Permissions, "Updated permissions mismatch")
}

func (s *roleTestSuite) testDelete(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)

	err := s.repo.Delete(ctx, role.ID)
	require.NoError(t, err, "Failed to delete role")

	_, err = s.repo.GetByID(ctx, role.ID)
	assert.Error(t, err, "Expected error when getting deleted role")
}

func (s *roleTestSuite) testList(t *testing.T) {
	ctx := context.Background()

	// Clear existing roles
	_, err := s.client.Role.Delete().Exec(ctx)
	require.NoError(t, err, "Failed to clear existing roles")

	for i := 0; i < 5; i++ {
		s.createRole(t)
	}

	roles, err := s.repo.List(ctx, 0, 10)
	require.NoError(t, err, "Failed to list roles")
	assert.Len(t, roles, 5, "Expected 5 roles")
}

func (s *roleTestSuite) testCount(t *testing.T) {
	ctx := context.Background()

	initialCount, err := s.repo.Count(ctx)
	require.NoError(t, err, "Failed to count roles")

	for i := 0; i < 5; i++ {
		s.createRole(t)
	}

	count, err := s.repo.Count(ctx)
	require.NoError(t, err, "Failed to count roles")
	assert.Equal(t, initialCount+5, count, "Expected count to increase by 5")
}

func (s *roleTestSuite) testAddPermission(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)

	err := s.repo.AddPermission(ctx, role.ID, "delete")
	require.NoError(t, err, "Failed to add permission")

	updatedRole, err := s.repo.GetByID(ctx, role.ID)
	require.NoError(t, err, "Failed to get updated role")
	assert.Contains(t, updatedRole.Permissions, "delete", "Added permission not found")
}

func (s *roleTestSuite) testRemovePermission(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)

	err := s.repo.RemovePermission(ctx, role.ID, "write")
	require.NoError(t, err, "Failed to remove permission")

	updatedRole, err := s.repo.GetByID(ctx, role.ID)
	require.NoError(t, err, "Failed to get updated role")
	assert.NotContains(t, updatedRole.Permissions, "write", "Removed permission still present")
}

func (s *roleTestSuite) testGetUsersInRole(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)
	user := s.createUser(t)

	err := s.repo.AddUserToRole(ctx, role.ID, user.ID)
	require.NoError(t, err, "Failed to add user to role")

	users, err := s.repo.GetUsersInRole(ctx, role.ID)
	require.NoError(t, err, "Failed to get users in role")
	assert.Len(t, users, 1, "Expected 1 user in role")
	assert.Equal(t, user.ID, users[0].ID, "User ID mismatch")
}

func (s *roleTestSuite) testSearch(t *testing.T) {
	ctx := context.Background()

	// Clear existing roles
	_, err := s.client.Role.Delete().Exec(ctx)
	require.NoError(t, err, "Failed to clear existing roles")

	// Create a role with a specific prefix for searching
	searchPrefix := "TestSearchRole_"
	roleName := searchPrefix + s.faker.UUID()
	role, err := s.repo.Create(ctx, &ent.Role{
		Name:        roleName,
		Permissions: []string{"read", "write"},
	})
	require.NoError(t, err, "Failed to create role for search test")

	// Create some additional roles to ensure our search is specific
	for i := 0; i < 5; i++ {
		s.createRole(t)
	}

	// Search for the specific role
	results, err := s.repo.Search(ctx, searchPrefix)
	require.NoError(t, err, "Failed to search roles")
	require.NotEmpty(t, results, "Expected search results")
	require.Len(t, results, 1, "Expected exactly one search result")
	assert.Equal(t, role.ID, results[0].ID, "Search result mismatch")
	assert.Equal(t, roleName, results[0].Name, "Role name mismatch in search result")
}

func (s *roleTestSuite) testGetRolesByUserID(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)
	user := s.createUser(t)

	err := s.repo.AssignRoleToUser(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to assign role to user")

	roles, err := s.repo.GetRolesByUserID(ctx, user.ID)
	require.NoError(t, err, "Failed to get roles by user ID")
	assert.Len(t, roles, 1, "Expected 1 role")
	assert.Equal(t, role.ID, roles[0].ID, "Role ID mismatch")
}

func (s *roleTestSuite) testAssignRoleToUser(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)
	user := s.createUser(t)

	err := s.repo.AssignRoleToUser(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to assign role to user")

	roles, err := s.repo.GetRolesByUserID(ctx, user.ID)
	require.NoError(t, err, "Failed to get roles by user ID")
	assert.Len(t, roles, 1, "Expected 1 role")
	assert.Equal(t, role.ID, roles[0].ID, "Role ID mismatch")
}

func (s *roleTestSuite) testRemoveRoleFromUser(t *testing.T) {
	ctx := context.Background()
	role := s.createRole(t)
	user := s.createUser(t)

	err := s.repo.AssignRoleToUser(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to assign role to user")

	err = s.repo.RemoveRoleFromUser(ctx, user.ID, role.ID)
	require.NoError(t, err, "Failed to remove role from user")

	roles, err := s.repo.GetRolesByUserID(ctx, user.ID)
	require.NoError(t, err, "Failed to get roles by user ID")
	assert.Empty(t, roles, "Expected no roles")
}

func (s *roleTestSuite) createUser(t *testing.T) *ent.User {
	user, err := s.client.User.Create().
		SetUsername(s.faker.Username()).
		SetEmail(s.faker.Email()).
		SetPassword(s.faker.Password(true, true, true, true, false, 32)).
		Save(context.Background())
	require.NoError(t, err, "Failed to create user")
	return user
}
