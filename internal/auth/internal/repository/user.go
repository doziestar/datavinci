// Package repository provides interfaces and implementations for user management.
package repository

import (
	"auth/ent"
	"auth/ent/role"
	"auth/ent/user"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// IUserRepository defines the interface for user-related operations.
type IUserRepository interface {
	// Create creates a new user in the database.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - u: Pointer to the ent.User object containing the user information to be created.
	//
	// Returns:
	//   - *ent.User: Pointer to the created user object.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   newUser := &ent.User{
	//     Username: "johndoe",
	//     Email:    "john@example.com",
	//     Password: "securepassword",
	//   }
	//   createdUser, err := repo.Create(ctx, newUser)
	//   if err != nil {
	//     log.Printf("Failed to create user: %v", err)
	//     return
	//   }
	//   fmt.Printf("Created user with ID: %s\n", createdUser.ID)
	Create(ctx context.Context, u *ent.User) (*ent.User, error)

	// GetByID retrieves a user by their ID.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - id: String representing the unique identifier of the user.
	//
	// Returns:
	//   - *ent.User: Pointer to the retrieved user object.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   userID := "123e4567-e89b-12d3-a456-426614174000"
	//   user, err := repo.GetByID(ctx, userID)
	//   if err != nil {
	//     log.Printf("Failed to get user: %v", err)
	//     return
	//   }
	//   fmt.Printf("Retrieved user: %s\n", user.Username)
	GetByID(ctx context.Context, id string) (*ent.User, error)

	// GetByUsername retrieves a user by their username.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - username: String representing the username of the user.
	//
	// Returns:
	//   - *ent.User: Pointer to the retrieved user object.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   username := "johndoe"
	//   user, err := repo.GetByUsername(ctx, username)
	//   if err != nil {
	//     log.Printf("Failed to get user: %v", err)
	//     return
	//   }
	//   fmt.Printf("Retrieved user: %s\n", user.Email)
	GetByUsername(ctx context.Context, username string) (*ent.User, error)

	// GetByEmail retrieves a user by their email address.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - email: String representing the email address of the user.
	//
	// Returns:
	//   - *ent.User: Pointer to the retrieved user object.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   email := "john@example.com"
	//   user, err := repo.GetByEmail(ctx, email)
	//   if err != nil {
	//     log.Printf("Failed to get user: %v", err)
	//     return
	//   }
	//   fmt.Printf("Retrieved user: %s\n", user.Username)
	GetByEmail(ctx context.Context, email string) (*ent.User, error)

	// Update updates an existing user's information.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - u: Pointer to the ent.User object containing the updated user information.
	//
	// Returns:
	//   - *ent.User: Pointer to the updated user object.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   user.Email = "newemail@example.com"
	//   updatedUser, err := repo.Update(ctx, user)
	//   if err != nil {
	//     log.Printf("Failed to update user: %v", err)
	//     return
	//   }
	//   fmt.Printf("Updated user email: %s\n", updatedUser.Email)
	Update(ctx context.Context, u *ent.User) (*ent.User, error)

	// Delete removes a user from the database by their ID.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - id: String representing the unique identifier of the user to be deleted.
	//
	// Returns:
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   userID := "123e4567-e89b-12d3-a456-426614174000"
	//   err := repo.Delete(ctx, userID)
	//   if err != nil {
	//     log.Printf("Failed to delete user: %v", err)
	//     return
	//   }
	//   fmt.Println("User successfully deleted")
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - offset: Integer representing the number of records to skip.
	//   - limit: Integer representing the maximum number of records to return.
	//
	// Returns:
	//   - []*ent.User: Slice of pointers to user objects.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   users, err := repo.List(ctx, 0, 10)
	//   if err != nil {
	//     log.Printf("Failed to list users: %v", err)
	//     return
	//   }
	//   for _, user := range users {
	//     fmt.Printf("User: %s, Email: %s\n", user.Username, user.Email)
	//   }
	List(ctx context.Context, offset, limit int) ([]*ent.User, error)

	// Count returns the total number of users in the database.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//
	// Returns:
	//   - int: The total number of users.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   count, err := repo.Count(ctx)
	//   if err != nil {
	//     log.Printf("Failed to count users: %v", err)
	//     return
	//   }
	//   fmt.Printf("Total number of users: %d\n", count)
	Count(ctx context.Context) (int, error)

	// AddRole assigns a role to a user.
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
	//   userID := "123e4567-e89b-12d3-a456-426614174000"
	//   roleID := "789e0123-e89b-12d3-a456-426614174000"
	//   err := repo.AddRole(ctx, userID, roleID)
	//   if err != nil {
	//     log.Printf("Failed to add role to user: %v", err)
	//     return
	//   }
	//   fmt.Println("Role successfully added to user")
	AddRole(ctx context.Context, userID, roleID string) error

	// RemoveRole removes a role from a user.
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
	//   userID := "123e4567-e89b-12d3-a456-426614174000"
	//   roleID := "789e0123-e89b-12d3-a456-426614174000"
	//   err := repo.RemoveRole(ctx, userID, roleID)
	//   if err != nil {
	//     log.Printf("Failed to remove role from user: %v", err)
	//     return
	//   }
	//   fmt.Println("Role successfully removed from user")
	RemoveRole(ctx context.Context, userID, roleID string) error

	// GetRoles retrieves all roles assigned to a user.
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
	//   userID := "123e4567-e89b-12d3-a456-426614174000"
	//   roles, err := repo.GetRoles(ctx, userID)
	//   if err != nil {
	//     log.Printf("Failed to get user roles: %v", err)
	//     return
	//   }
	//   for _, role := range roles {
	//     fmt.Printf("Role: %s\n", role.Name)
	//   }
	GetRoles(ctx context.Context, userID string) ([]*ent.Role, error)

	// Search performs a search for users based on a query string.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - query: String representing the search query.
	//
	// Returns:
	//   - []*ent.User: Slice of pointers to user objects matching the search query.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   query := "john"
	//   users, err := repo.Search(ctx, query)
	//   if err != nil {
	//     log.Printf("Failed to search users: %v", err)
	//     return
	//   }
	//   for _, user := range users {
	//     fmt.Printf("Matching user: %s, Email: %s\n", user.Username, user.Email)
	//   }
	Search(ctx context.Context, query string) ([]*ent.User, error)

	// ChangePassword updates a user's password.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - userID: String representing the unique identifier of the user.
	//   - newPassword: String representing the new password.
	//
	// Returns:
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   userID := "123e4567-e89b-12d3-a456-426614174000"
	//   newPassword := "newSecurePassword123"
	//   err := repo.ChangePassword(ctx, userID, newPassword)
	//   if err != nil {
	//     log.Printf("Failed to change password: %v", err)
	//     return
	//   }
	//   fmt.Println("Password successfully changed")
	ChangePassword(ctx context.Context, userID, newPassword string) error

	// GetUsersByRole retrieves all users with a specific role.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - roleID: String representing the unique identifier of the role.
	//
	// Returns:
	//   - []*ent.User: Slice of pointers to user objects with the specified role.
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   roleID := "789e0123-e89b-12d3-a456-426614174000"
	//   users, err := repo.GetUsersByRole(ctx, roleID)
	//   if err != nil {
	//     log.Printf("Failed to get users by role: %v", err)
	//     return
	//   }
	//   for _, user := range users {
	//     fmt.Printf("User with role: %s, Email: %s\n", user.Username, user.Email)
	//   }
	GetUsersByRole(ctx context.Context, roleID string) ([]*ent.User, error)

	// CheckPassword verifies if the provided password is correct for a user.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - password: String representing the password to check.
	//
	// Returns:
	//   - bool: True if the password is correct, false otherwise.
	//
	// Example:
	//   password := "userPassword123"
	//   isCorrect := repo.CheckPassword(ctx, password)
	//   if isCorrect {
	//     fmt.Println("Password is correct")
	//   } else {
	//     fmt.Println("Password is incorrect")
	//   }
	CheckPassword(ctx context.Context, password string) bool

	// SetPassword sets a new password for a user.
	//
	// Parameters:
	//   - ctx: Context for the database operation.
	//   - password: String representing the new password to set.
	//
	// Returns:
	//   - error: An error object that indicates success or failure of the operation.
	//
	// Example:
	//   newPassword := "newSecurePassword123"
	//   err := repo.SetPassword(ctx, newPassword)
	//   if err != nil {
	//     log.Printf("Failed to set new password: %v", err)
	//     return
	//   }
	//   fmt.Println("New password successfully set")
	SetPassword(ctx context.Context, password string) error
}

type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new instance of UserRepository with the given ent.Client.
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	return r.client.User.
		Create().
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPassword(u.Password).
		Save(ctx)
}

// GetByID retrieves a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*ent.User, error) {
	return r.client.User.Query().Where(user.ID(id)).Only(ctx)
}

// GetByUsername retrieves a user by their username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(user.Username(username)).Only(ctx)
}

// GetByEmail retrieves a user by their email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().Where(user.Email(email)).Only(ctx)
}

// Update updates an existing user's information.
func (r *UserRepository) Update(ctx context.Context, u *ent.User) (*ent.User, error) {
	return r.client.User.UpdateOne(u).
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

// Delete removes a user from the database by their ID.
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*ent.User, error) {
	return r.client.User.Query().
		Offset(offset).
		Limit(limit).
		All(ctx)
}

// Count returns the total number of users in the database.
func (r *UserRepository) Count(ctx context.Context) (int, error) {
	return r.client.User.Query().Count(ctx)
}

// AddRole assigns a role to a user.
func (r *UserRepository) AddRole(ctx context.Context, userID, roleID string) error {
	return r.client.User.UpdateOneID(userID).
		AddRoleIDs(roleID).
		Exec(ctx)
}

// RemoveRole removes a role from a user.
func (r *UserRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	return r.client.User.UpdateOneID(userID).
		RemoveRoleIDs(roleID).
		Exec(ctx)
}

// GetRoles retrieves all roles assigned to a user.
func (r *UserRepository) GetRoles(ctx context.Context, userID string) ([]*ent.Role, error) {
	u, err := r.client.User.Query().
		Where(user.ID(userID)).
		WithRoles().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return u.Edges.Roles, nil
}

// Search performs a search for users based on a query string.
func (r *UserRepository) Search(ctx context.Context, query string) ([]*ent.User, error) {
	return r.client.User.Query().
		Where(
			user.Or(
				user.UsernameContains(query),
				user.EmailContains(query),
			),
		).
		All(ctx)
}

// ChangePassword updates a user's password.
func (r *UserRepository) ChangePassword(ctx context.Context, userID, newPassword string) error {
	u, err := r.client.User.Query().Where(user.ID(userID)).Only(ctx)
	if err != nil {
		return err
	}

	return r.client.User.UpdateOne(u).
		SetPassword(newPassword).
		SetUpdatedAt(time.Now()).
		Exec(ctx)

}

// GetUsersByRole retrieves all users with a specific role.
func (r *UserRepository) GetUsersByRole(ctx context.Context, roleID string) ([]*ent.User, error) {
	return r.client.User.Query().
		Where(user.HasRolesWith(role.ID(roleID))).
		All(ctx)
}

// CheckPassword verifies if the provided password is correct for a user.
func (r *UserRepository) CheckPassword(ctx context.Context, password string) bool {
	user, err := r.client.User.Query().Where(user.Password(password)).Only(ctx)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// SetPassword sets a new password for a user.
func (r *UserRepository) SetPassword(ctx context.Context, password string) error {
	user, err := r.client.User.Query().Where(user.Password(password)).Only(ctx)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return r.client.User.UpdateOne(user).SetPassword(string(hashedPassword)).Exec(ctx)
}
