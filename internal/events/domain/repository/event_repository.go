package repository

import (
	"context"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
)

type EventRepository interface {
	Save(ctx context.Context, event *entity.Event) (*entity.Event, error)
	FindByID(ctx context.Context, id string) (*entity.Event, error)
	FindAll(ctx context.Context) ([]*entity.Event, error)
	Update(ctx context.Context, event *entity.Event) (*entity.Event, error)
	PartialUpdate(ctx context.Context, id string, input UpdateEventInput) (*entity.Event, error)
	Delete(ctx context.Context, id string) error
	FindByIDWithActivitiesAndCheckIns(ctx context.Context, eventID string) (*EventWithActivitiesAndCheckIns, error)
}

// UpdateEventInput representa os campos opcionais para atualização parcial
type UpdateEventInput struct {
	Name           *string
	AllowedDomains *[]string
	Description    *string
	StartDate      *time.Time
	EndDate        *time.Time
	Status         *entity.EventStatus
}

// Query Results

type EventWithActivitiesAndCheckIns struct {
	Event      *entity.Event
	Activities []ActivityWithCheckIns
}

type ActivityWithCheckIns struct {
	ActivityID   string           `json:"activity_id"`
	ActivityName string           `json:"activity_name"`
	CheckIns     []entity.CheckIn `json:"check_ins"`
}
