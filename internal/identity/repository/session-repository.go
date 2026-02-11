package identity

import (
	"context"

	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/domain"
)



type SessionRepository interface {
	Save(ctx context.Context, session *identity.Session) error
	FindByRefreshToken(ctx context.Context, token string) (*identity.Session, error)
	FindByUserID(ctx context.Context, userID string) ([]*identity.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}