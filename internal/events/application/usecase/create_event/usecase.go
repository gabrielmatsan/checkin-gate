package createevent

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/service"
)

type Input struct {
	UserID         string
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
	eventRepo   repository.EventRepository
	userAuthSvc service.UserAuthorizationService
}

func NewUseCase(eventRepo repository.EventRepository, userAuthSvc service.UserAuthorizationService) *UseCase {
	return &UseCase{
		eventRepo:   eventRepo,
		userAuthSvc: userAuthSvc,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	user, err := uc.userAuthSvc.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if user == nil {
		return nil, service.ErrUserNotFound
	}

	if !user.IsAdmin {
		return nil, service.ErrUserNotAdmin
	}

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

	newEvent, err := uc.eventRepo.Save(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("failed to save event: %w", err)
	}

	return &Output{
		Event: newEvent,
	}, nil
}
