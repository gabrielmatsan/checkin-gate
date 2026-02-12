package events

import (
	"context"
	"fmt"

	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
	repository "github.com/gabrielmatsan/checkin-gate/internal/events/repository"
)

type GetEventWithActivitiesInput struct {
	EventID string `json:"event_id"`
}

type EventWithActivities struct {
	Event      *events.Event      `json:"event"`
	Activities []*events.Activity `json:"activities"`
}

type GetEventWithActivitiesUseCase struct {
	eventRepo    repository.EventRepository
	activityRepo repository.ActivityRepository
}

func NewGetEventWithActivitiesUseCase(eventRepo repository.EventRepository, activityRepo repository.ActivityRepository) *GetEventWithActivitiesUseCase {
	return &GetEventWithActivitiesUseCase{
		eventRepo:    eventRepo,
		activityRepo: activityRepo,
	}
}

func (uc *GetEventWithActivitiesUseCase) Execute(ctx context.Context, input GetEventWithActivitiesInput) (*EventWithActivities, error) {

	event, err := uc.eventRepo.FindByID(ctx, input.EventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, fmt.Errorf("event with ID %s not found", input.EventID)
	}

	activities, err := uc.activityRepo.FindByEventID(ctx, input.EventID)
	if err != nil {
		return nil, err
	}

	return &EventWithActivities{
		Event:      event,
		Activities: activities,
	}, nil
}
