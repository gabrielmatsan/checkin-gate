package events

import (
	"context"

	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
)


type EventRepository interface {
	Save(ctx context.Context, event *events.Event) (*events.Event, error)
	FindByID(ctx context.Context, id string) (*events.Event, error)
	FindAll(ctx context.Context) ([]*events.Event, error)
	Update(ctx context.Context, event *events.Event) (*events.Event, error)
	Delete(ctx context.Context, id string) error
}
