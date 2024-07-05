package repository

import (
	"auth/ent"
	"context"
	"time"

	"auth/ent/role"
	"auth/ent/user"
)

// IRoleRepository defines the interface for role-related operations.
type IRoleRepository interface {
    // Create creates a new role in the database.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - role: Pointer to the ent.Role object containing the role information to be created.
    //
    // Returns:
    //   - *ent.Role: Pointer to the created role object.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   newRole := &ent.Role{
    //     Name: "admin",
    //     Description: "Administrator role with full access",
    //   }
    //   createdRole, err := repo.Create(ctx, newRole)
    //   if err != nil {
    //     log.Printf("Failed to create role: %v", err)
    //     return
    //   }
    //   fmt.Printf("Created role with ID: %s\n", createdRole.ID)
    Create(ctx context.Context, role *ent.Role) (*ent.Role, error)

    // GetByID retrieves a role by its ID.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - id: String representing the unique identifier of the role.
    //
    // Returns:
    //   - *ent.Role: Pointer to the retrieved role object.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   role, err := repo.GetByID(ctx, roleID)
    //   if err != nil {
    //     log.Printf("Failed to get role: %v", err)
    //     return
    //   }
    //   fmt.Printf("Retrieved role: %s\n", role.Name)
    GetByID(ctx context.Context, id string) (*ent.Role, error)

    // GetByName retrieves a role by its name.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - name: String representing the name of the role.
    //
    // Returns:
    //   - *ent.Role: Pointer to the retrieved role object.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleName := "admin"
    //   role, err := repo.GetByName(ctx, roleName)
    //   if err != nil {
    //     log.Printf("Failed to get role: %v", err)
    //     return
    //   }
    //   fmt.Printf("Retrieved role ID: %s\n", role.ID)
    GetByName(ctx context.Context, name string) (*ent.Role, error)

    // Update updates an existing role's information.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - role: Pointer to the ent.Role object containing the updated role information.
    //
    // Returns:
    //   - *ent.Role: Pointer to the updated role object.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   role.Description = "Updated administrator role description"
    //   updatedRole, err := repo.Update(ctx, role)
    //   if err != nil {
    //     log.Printf("Failed to update role: %v", err)
    //     return
    //   }
    //   fmt.Printf("Updated role description: %s\n", updatedRole.Description)
    Update(ctx context.Context, role *ent.Role) (*ent.Role, error)

    // Delete removes a role from the database by its ID.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - id: String representing the unique identifier of the role to be deleted.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   err := repo.Delete(ctx, roleID)
    //   if err != nil {
    //     log.Printf("Failed to delete role: %v", err)
    //     return
    //   }
    //   fmt.Println("Role successfully deleted")
    Delete(ctx context.Context, id string) error

    // List retrieves a paginated list of roles.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - offset: Integer representing the number of records to skip.
    //   - limit: Integer representing the maximum number of records to return.
    //
    // Returns:
    //   - []*ent.Role: Slice of pointers to role objects.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roles, err := repo.List(ctx, 0, 10)
    //   if err != nil {
    //     log.Printf("Failed to list roles: %v", err)
    //     return
    //   }
    //   for _, role := range roles {
    //     fmt.Printf("Role: %s, Description: %s\n", role.Name, role.Description)
    //   }
    List(ctx context.Context, offset, limit int) ([]*ent.Role, error)

    // Count returns the total number of roles in the database.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //
    // Returns:
    //   - int: The total number of roles.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   count, err := repo.Count(ctx)
    //   if err != nil {
    //     log.Printf("Failed to count roles: %v", err)
    //     return
    //   }
    //   fmt.Printf("Total number of roles: %d\n", count)
    Count(ctx context.Context) (int, error)

    // AddPermission adds a permission to a role.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - roleID: String representing the unique identifier of the role.
    //   - permission: String representing the permission to be added.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   permission := "create:user"
    //   err := repo.AddPermission(ctx, roleID, permission)
    //   if err != nil {
    //     log.Printf("Failed to add permission to role: %v", err)
    //     return
    //   }
    //   fmt.Println("Permission successfully added to role")
    AddPermission(ctx context.Context, roleID, permission string) error

    // RemovePermission removes a permission from a role.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - roleID: String representing the unique identifier of the role.
    //   - permission: String representing the permission to be removed.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   permission := "delete:user"
    //   err := repo.RemovePermission(ctx, roleID, permission)
    //   if err != nil {
    //     log.Printf("Failed to remove permission from role: %v", err)
    //     return
    //   }
    //   fmt.Println("Permission successfully removed from role")
    RemovePermission(ctx context.Context, roleID, permission string) error

    // GetPermissions retrieves all permissions assigned to a role.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - roleID: String representing the unique identifier of the role.
    //
    // Returns:
    //   - []string: Slice of strings representing the permissions assigned to the role.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   permissions, err := repo.GetPermissions(ctx, roleID)
    //   if err != nil {
    //     log.Printf("Failed to get role permissions: %v", err)
    //     return
    //   }
    //   for _, perm := range permissions {
    //     fmt.Printf("Permission: %s\n", perm)
    //   }
    GetPermissions(ctx context.Context, roleID string) ([]string, error)

    // AddUserToRole assigns a user to a role.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - roleID: String representing the unique identifier of the role.
    //   - userID: String representing the unique identifier of the user to be added to the role.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   userID := "456e7890-e89b-12d3-a456-426614174000"
    //   err := repo.AddUserToRole(ctx, roleID, userID)
    //   if err != nil {
    //     log.Printf("Failed to add user to role: %v", err)
    //     return
    //   }
    //   fmt.Println("User successfully added to role")
    AddUserToRole(ctx context.Context, roleID, userID string) error

    // RemoveUserFromRole removes a user from a role.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - roleID: String representing the unique identifier of the role.
    //   - userID: String representing the unique identifier of the user to be removed from the role.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   userID := "456e7890-e89b-12d3-a456-426614174000"
    //   err := repo.RemoveUserFromRole(ctx, roleID, userID)
    //   if err != nil {
    //     log.Printf("Failed to remove user from role: %v", err)
    //     return
    //   }
    //   fmt.Println("User successfully removed from role")
    RemoveUserFromRole(ctx context.Context, roleID, userID string) error

    // GetUsersInRole retrieves all users assigned to a specific role.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - roleID: String representing the unique identifier of the role.
    //
    // Returns:
    //   - []*ent.User: Slice of pointers to user objects assigned to the role.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   users, err := repo.GetUsersInRole(ctx, roleID)
    //   if err != nil {
    //     log.Printf("Failed to get users in role: %v", err)
    //     return
    //   }
    //   for _, user := range users {
    //     fmt.Printf("User in role: %s, Email: %s\n", user.Username, user.Email)
    //   }
    GetUsersInRole(ctx context.Context, roleID string) ([]*ent.User, error)

    // Search performs a search for roles based on a query string.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - query: String representing the search query.
    //
    // Returns:
    //   - []*ent.Role: Slice of pointers to role objects matching the search query.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   query := "admin"
    //   roles, err := repo.Search(ctx, query)
    //   if err != nil {
    //     log.Printf("Failed to search roles: %v", err)
    //     return
    //   }
    //   for _, role := range roles {
    //     fmt.Printf("Matching role: %s, Description: %s\n", role.Name, role.Description)
    //   }
    Search(ctx context.Context, query string) ([]*ent.Role, error)

    // GetRolesByUserID retrieves all roles assigned to a specific user.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - userID: String representing the unique identifier of the user.
    //
    // Returns:
    //   - []*ent.Role: Slice of pointers to role objects assigned to the user.
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   userID := "456e7890-e89b-12d3-a456-426614174000"
    //   roles, err := repo.GetRolesByUserID(ctx, userID)
    //   if err != nil {
    //     log.Printf("Failed to get roles for user: %v", err)
    //     return
    //   }
    //   for _, role := range roles {
    //     fmt.Printf("User's role: %s, Description: %s\n", role.Name, role.Description)
    //   }
    GetRolesByUserID(ctx context.Context, userID string) ([]*ent.Role, error)

    // AssignRoleToUser assigns a role to a user.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - userID: String representing the unique identifier of the user.
    //   - roleID: String representing the unique identifier of the role to be assigned.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   userID := "456e7890-e89b-12d3-a456-426614174000"
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   err := repo.AssignRoleToUser(ctx, userID, roleID)
    //   if err != nil {
    //     log.Printf("Failed to assign role to user: %v", err)
    //     return
    //   }
    //   fmt.Println("Role successfully assigned to user")
    AssignRoleToUser(ctx context.Context, userID, roleID string) error

    // RemoveRoleFromUser removes a role from a user.
    //
    // Parameters:
    //   - ctx: Context for the database operation.
    //   - userID: String representing the unique identifier of the user.
    //   - roleID: String representing the unique identifier of the role to be removed.
    //
    // Returns:
    //   - error: An error object that indicates success or failure of the operation.
    //
    // Example:
    //   userID := "456e7890-e89b-12d3-a456-426614174000"
    //   roleID := "123e4567-e89b-12d3-a456-426614174000"
    //   err := repo.RemoveRoleFromUser(ctx, userID, roleID)
    //   if err != nil {
    //     log.Printf("Failed to remove role from user: %v", err)
    //     return
    //   }
    //   fmt.Println("Role successfully removed from user")
    RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
}

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
    roleData, err := r.GetByID(ctx, roleID)
    if err != nil {
        return err
    }

    updatedPermissions := make([]string, 0)
    for _, p := range roleData.Permissions {
        if p != permission {
            updatedPermissions = append(updatedPermissions, p)
        }
    }

    var newName = updatedPermissions

    println(newName)

    return r.client.Role.UpdateOneID(roleID).
        SetPermissions(updatedPermissions).
        Exec(ctx)
}

func (r *RoleRepository) GetPermissions(ctx context.Context, roleID string) ([]string, error) {
    roleData, err := r.GetByID(ctx, roleID)
    if err != nil {
        return nil, err
    }
    return roleData.Permissions, nil
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
    roleData, err := r.client.Role.Query().
        Where(role.ID(roleID)).
        WithUsers().
        Only(ctx)
    if err != nil {
        return nil, err
    }
    return roleData.Edges.Users, nil
}

func (r *RoleRepository) Search(ctx context.Context, query string) ([]*ent.Role, error) {
    return r.client.Role.Query().
        Where(role.NameHasPrefix(query)).
        All(ctx)
}

func (r *RoleRepository) GetRolesByUserID(ctx context.Context, userID string) ([]*ent.Role, error) {
    return r.client.Role.Query().
        Where(role.HasUsersWith(user.ID(userID))).
        All(ctx)
}

func (r *RoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
    return r.client.User.UpdateOneID(userID).
        AddRoleIDs(roleID).
        Exec(ctx)
}

func (r *RoleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
    return r.client.User.UpdateOneID(userID).
        RemoveRoleIDs(roleID).
        Exec(ctx)
}

// func (r *RoleRepository) GetRolesByPermission(ctx context.Context, permission string) ([]*ent.Role, error) {
//     return r.client.Role.Query().
//         Where(role.HasPermissionsWith(permission)).
//         All(ctx)
// }
