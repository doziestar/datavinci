package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return uuid.New().String()
			}),
		field.String("username").Unique(),
		field.String("email").Unique(),
		field.String("password"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", Role.Type),
		edge.To("tokens", Token.Type),
	}
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
    return []ent.Hook{
        HashPassword(),
    }
}

// HashPassword is a hook that hashes the password before creating or updating a user.
func HashPassword() ent.Hook {
    return func(next ent.Mutator) ent.Mutator {
        return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
            if password, ok := m.Field("password"); ok {
                hash, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
                if err != nil {
                    return nil, err
                }
                m.SetField("password", string(hash))
            }
            return next.Mutate(ctx, m)
        })
    }
}