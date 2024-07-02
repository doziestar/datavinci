package repository

import (
	"auth/ent"
	"context"
	"time"

	"auth/ent/role"
)

type RoleRepository struct {
	client *ent.Client
}

func NewRoleRepository(client *ent.Client) *RoleRepository {
	return &RoleRepository{client: client}
}

func (r *RoleRepository) Create(ctx context.Context, role *ent.Role) (*ent.Role, error) {
	return r.client.Role.
		Create().
		SetName(role.Name).
		SetPermissions(role.Permissions).
		Save(ctx)
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*ent.Role, error) {
	return r.client.Role.Query().Where(role.ID(id)).Only(ctx)
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*ent.Role, error) {
	return r.client.Role.Query().Where(role.Name(name)).Only(ctx)
}

func (r *RoleRepository) Update(ctx context.Context, role *ent.Role) (*ent.Role, error) {
	return r.client.Role.UpdateOne(role).
		SetName(role.Name).
		SetPermissions(role.Permissions).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	return r.client.Role.DeleteOneID(id).Exec(ctx)
}

func (r *RoleRepository) List(ctx context.Context, offset, limit int) ([]*ent.Role, error) {
	return r.client.Role.Query().
		Offset(offset).
		Limit(limit).
		All(ctx)
}

func (r *RoleRepository) Count(ctx context.Context) (int, error) {
	return r.client.Role.Query().Count(ctx)
}

func (r *RoleRepository) AddPermission(ctx context.Context, roleID, permission string) error {
	return r.client.Role.UpdateOneID(roleID).
		AppendPermissions([]string{permission}).
		Exec(ctx)
}

func (r *RoleRepository) RemovePermission(ctx context.Context, roleID, permission string) error {
	role, err := r.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	updatedPermissions := make([]string, 0)
	for _, p := range role.Permissions {
		if p != permission {
			updatedPermissions = append(updatedPermissions, p)
		}
	}

	return r.client.Role.UpdateOneID(roleID).
		SetPermissions(updatedPermissions).
		Exec(ctx)
}

func (r *RoleRepository) GetPermissions(ctx context.Context, roleID string) ([]string, error) {
	role, err := r.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (r *RoleRepository) AddUserToRole(ctx context.Context, roleID, userID string) error {
	return r.client.Role.UpdateOneID(roleID).
		AddUserIDs(userID).
		Exec(ctx)
}

func (r *RoleRepository) RemoveUserFromRole(ctx context.Context, roleID, userID string) error {
	return r.client.Role.UpdateOneID(roleID).
		RemoveUserIDs(userID).
		Exec(ctx)
}

func (r *RoleRepository) GetUsersInRole(ctx context.Context, roleID string) ([]*ent.User, error) {
	role, err := r.client.Role.Query().
		Where(role.ID(roleID)).
		WithUsers().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return role.Edges.Users, nil
}

func (r *RoleRepository) Search(ctx context.Context, query string) ([]*ent.Role, error) {
	return r.client.Role.Query().
		Where(role.NameContains(query)).
		All(ctx)
}

// func (r *RoleRepository) GetRolesByPermission(ctx context.Context, permission string) ([]*ent.Role, error) {
//     return r.client.Role.Query().
//         Where(role.HasPermissionsWith(permission)).
//         All(ctx)
// }
