package repository

import (
	"auth/ent"
	"context"
	"time"

	"auth/ent/role"
	"auth/ent/user"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	client *ent.Client
}

// NewUserRepository creates a new instance of UserRepository with the given ent.Client.
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, u *ent.User) (*ent.User, error) {
	return r.client.User.
		Create().
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPassword(u.Password).
		Save(ctx)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*ent.User, error) {
	return r.client.User.Query().Where(user.ID(id)).Only(ctx)
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().Where(user.Username(username)).Only(ctx)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().Where(user.Email(email)).Only(ctx)
}

func (r *UserRepository) Update(ctx context.Context, u *ent.User) (*ent.User, error) {
	return r.client.User.UpdateOne(u).
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetUpdatedAt(time.Now()).
		Save(ctx)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*ent.User, error) {
	return r.client.User.Query().
		Offset(offset).
		Limit(limit).
		All(ctx)
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	return r.client.User.Query().Count(ctx)
}

func (r *UserRepository) AddRole(ctx context.Context, userID, roleID string) error {
	return r.client.User.UpdateOneID(userID).
		AddRoleIDs(roleID).
		Exec(ctx)
}

func (r *UserRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	return r.client.User.UpdateOneID(userID).
		RemoveRoleIDs(roleID).
		Exec(ctx)
}

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

func (r *UserRepository) GetUsersByRole(ctx context.Context, roleID string) ([]*ent.User, error) {
	return r.client.User.Query().
		Where(user.HasRolesWith(role.ID(roleID))).
		All(ctx)
}

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
