package events

import (
	"context"

	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
)

type ActivityRepository interface {
	Save(ctx context.Context, activity *events.Activity) (*events.Activity, error)
	SaveAll(ctx context.Context, activities []*events.Activity) ([]*events.Activity, error)
	FindByID(ctx context.Context, id string) (*events.Activity, error)
	FindByEventID(ctx context.Context, eventID string) ([]*events.Activity, error)
	FindByEventIDAndNames(ctx context.Context, eventID string, names []string) ([]*events.Activity, error)
	FindAll(ctx context.Context) ([]*events.Activity, error)
	Update(ctx context.Context, activity *events.Activity) (*events.Activity, error)
	Delete(ctx context.Context, id string) error
}
