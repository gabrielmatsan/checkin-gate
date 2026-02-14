package repository

import (
	"context"

	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
)

type SessionRepository interface {
	Save(ctx context.Context, session *entity.Session) error
	FindByRefreshToken(ctx context.Context, token string) (*entity.Session, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
