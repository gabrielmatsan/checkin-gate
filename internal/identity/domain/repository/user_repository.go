package repository

import (
	"context"

	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
