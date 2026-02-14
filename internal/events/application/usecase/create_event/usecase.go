package createevent

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
)

type Input struct {
	Name           string
	AllowedDomains *[]string
	Description    *string
	StartDate      time.Time
	EndDate        time.Time
}

type Output struct {
	Event *entity.Event
}

type UseCase struct {
	eventRepo repository.EventRepository
}

func NewUseCase(eventRepo repository.EventRepository) *UseCase {
	return &UseCase{
		eventRepo: eventRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	event, err := entity.NewEvent(entity.NewEventParams{
		Name:           input.Name,
		AllowedDomains: input.AllowedDomains,
		Description:    input.Description,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
	})
	if err != nil {
		return nil, err
	}

	if !event.IsStartDateBeforeEndDate() {
		return nil, fmt.Errorf("start date must be before end date")
	}

	if !event.IsEndDateAfterStartDate() {
		return nil, fmt.Errorf("end date must be after start date")
	}

	newEvent, err := uc.eventRepo.Save(ctx, event)
	if err != nil {
		return nil, err
	}

	return &Output{
		Event: newEvent,
	}, nil
}
