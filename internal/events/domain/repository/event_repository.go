package repository

import (
	"context"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
)

type EventRepository interface {
	Save(ctx context.Context, event *entity.Event) (*entity.Event, error)
	FindByID(ctx context.Context, id string) (*entity.Event, error)
	FindAll(ctx context.Context) ([]*entity.Event, error)
	Update(ctx context.Context, event *entity.Event) (*entity.Event, error)
	Delete(ctx context.Context, id string) error
}
