package repository

import (
	"context"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
)

// retornar evento com as atividades
type CheckInRepository interface {
	Save(ctx context.Context, checkIn *entity.CheckIn) (*entity.CheckIn, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.CheckIn, error)
	FindByActivityID(ctx context.Context, activityID string) ([]*entity.CheckIn, error)
	FindByID(ctx context.Context, id string) (*entity.CheckIn, error)
	FindByUserAndActivity(ctx context.Context, userID, activityID string) (*entity.CheckIn, error)
}
