package geteventwithactivities

import (
	"context"
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
)

type Input struct {
	EventID string
}

type Output struct {
	Event      *entity.Event
	Activities []*entity.Activity
}

type UseCase struct {
	eventRepo    repository.EventRepository
	activityRepo repository.ActivityRepository
}

func NewUseCase(
	eventRepo repository.EventRepository,
	activityRepo repository.ActivityRepository,
) *UseCase {
	return &UseCase{
		eventRepo:    eventRepo,
		activityRepo: activityRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
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

	return &Output{
		Event:      event,
		Activities: activities,
	}, nil
}
