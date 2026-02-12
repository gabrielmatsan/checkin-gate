package events

import (
	"context"
	"fmt"
	"time"

	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
	repository "github.com/gabrielmatsan/checkin-gate/internal/events/repository"
)

type CreateEventUseCase struct {
	eventRepo repository.EventRepository
}

type CreateEventInput struct {
	Name           string    `json:"name"`
	AllowedDomains *[]string `json:"allowed_domains"`
	Description    *string   `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

func NewCreateEventUseCase(eventRepo repository.EventRepository) *CreateEventUseCase {
	return &CreateEventUseCase{
		eventRepo: eventRepo,
	}
}

func (uc *CreateEventUseCase) Execute(ctx context.Context, input CreateEventInput) (*events.Event, error) {

	event, err := events.NewEvent(events.NewEventParams{
		Name:           input.Name,
		AllowedDomains: input.AllowedDomains,
		Description:    input.Description,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
	})

	if err != nil {
		return nil, err
	}

	// validar se startDate é menor que endDate e se endDate é maior que startDate
	if event.StartDate.After(event.EndDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	if event.EndDate.Before(event.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	newEvent, err := uc.eventRepo.Save(ctx, event)
	if err != nil {
		return nil, err
	}

	return newEvent, nil
}
