package repository

import (
	"context"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
)

type ActivityRepository interface {
	Save(ctx context.Context, activity *entity.Activity) (*entity.Activity, error)
	SaveAll(ctx context.Context, activities []*entity.Activity) ([]*entity.Activity, error)
	FindByID(ctx context.Context, id string) (*entity.Activity, error)
	FindByEventID(ctx context.Context, eventID string) ([]*entity.Activity, error)
	FindByEventIDAndNames(ctx context.Context, eventID string, names []string) ([]*entity.Activity, error)
	FindAll(ctx context.Context) ([]*entity.Activity, error)
	Update(ctx context.Context, activity *entity.Activity) (*entity.Activity, error)
	Delete(ctx context.Context, id string) error
}
