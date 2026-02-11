package identity

import (
	"context"

	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *identity.User) (*identity.User, error)
	FindByID(ctx context.Context, id string) (*identity.User, error)
	FindByEmail(ctx context.Context, email string) (*identity.User, error)
	Update(ctx context.Context, user *identity.User) error
	Delete(ctx context.Context, id string) error
}
