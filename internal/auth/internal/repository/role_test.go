package repository

import (
	"auth/ent"
	"auth/ent/enttest"
	_ "auth/ent/role"
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v7"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func createTestClient(t *testing.T) *ent.Client {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	t.Cleanup(func() {
		client.Close()
	})
	return client
}

func TestRoleRepositoryCreate(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := &ent.Role{
		Name:        fake.Word(),
		Permissions: []string{"read", "write"},
	}

	createdRole, err := repo.Create(ctx, newRole)
	assert.NoError(t, err)
	assert.NotNil(t, createdRole)
	assert.Equal(t, newRole.Name, createdRole.Name)
	assert.ElementsMatch(t, newRole.Permissions, createdRole.Permissions)
}

func TestRoleRepositoryGetByID(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName(fake.Name()).
		SetPermissions([]string{"read", "write"}).
		SaveX(ctx)

	fetchedRole, err := repo.GetByID(ctx, newRole.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedRole)
	assert.Equal(t, newRole.Name, fetchedRole.Name)
	assert.ElementsMatch(t, newRole.Permissions, fetchedRole.Permissions)
}

func TestRoleRepositoryGetByName(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName(fake.Name()).
		SetPermissions([]string{"read", "write"}).
		SaveX(ctx)

	fetchedRole, err := repo.GetByName(ctx, newRole.Name)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedRole)
	assert.Equal(t, newRole.Name, fetchedRole.Name)
	assert.ElementsMatch(t, newRole.Permissions, fetchedRole.Permissions)
}

func TestRoleRepositoryUpdate(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName(fake.Name()).
		SetPermissions([]string{"read"}).
		SaveX(ctx)

	newRole.Permissions = []string{"read", "write"}
	updatedRole, err := repo.Update(ctx, newRole)
	assert.NoError(t, err)
	assert.NotNil(t, updatedRole)
	assert.ElementsMatch(t, newRole.Permissions, updatedRole.Permissions)
}

func TestRoleRepositoryDelete(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName("test-role").
		SetPermissions([]string{"read"}).
		SaveX(ctx)

	err := repo.Delete(ctx, newRole.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, newRole.ID)
	assert.Error(t, err)
}

func TestRoleRepositoryList(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		client.Role.
			Create().
			SetName(fake.Name()).
			SetPermissions([]string{"read"}).
			SaveX(ctx)
	}

	roles, err := repo.List(ctx, 0, 10)
	assert.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 5)
}

func TestRoleRepositoryCount(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	for i := 0; i < 5; i++ {
		client.Role.
			Create().
			SetName(fake.Name()).
			SetPermissions([]string{"read"}).
			SaveX(ctx)
	}

	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestRoleRepositoryAddPermission(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName(fake.Name()).
		SetPermissions([]string{"read"}).
		SaveX(ctx)

	err := repo.AddPermission(ctx, newRole.ID, "write")
	assert.NoError(t, err)

	fetchedRole, err := repo.GetByID(ctx, newRole.ID)
	assert.NoError(t, err)
	assert.Contains(t, fetchedRole.Permissions, "write")
}

func TestRoleRepositoryRemovePermission(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName(fake.Name()).
		SetPermissions([]string{"read", "write"}).
		SaveX(ctx)

	err := repo.RemovePermission(ctx, newRole.ID, "write")
	assert.NoError(t, err)

	fetchedRole, err := repo.GetByID(ctx, newRole.ID)
	assert.NoError(t, err)
	assert.NotContains(t, fetchedRole.Permissions, "write")
}

func TestRoleRepositoryGetUsersInRole(t *testing.T) {
	client := createTestClient(t)
	repo := NewRoleRepository(client)

	ctx := context.Background()
	newRole := client.Role.
		Create().
		SetName(fake.Name()).
		SetPermissions([]string{"read"}).
		SaveX(ctx)
	newUser := client.User.
		Create().
		SetUsername(fake.Name()).
		SetEmail(fake.Email()).
		SetPassword(fake.Password(true, true, true, true, false, 12)).
		AddRoleIDs(newRole.ID).
		SaveX(ctx)

	users, err := repo.GetUsersInRole(ctx, newRole.ID)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 1)
	assert.Equal(t, newUser.ID, users[0].ID)
}
