package geteventdetails

import (
	"context"
	"fmt"
	"slices"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
)

type Input struct {
	EventID string
	Role    string
}

type Output struct {
	*repository.EventWithActivitiesAndCheckIns
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

	if !slices.Contains([]string{"admin", "user"}, input.Role) {
		return nil, fmt.Errorf("user is not authorized to get event details")
	}

	result, err := uc.eventRepo.FindByIDWithActivitiesAndCheckIns(ctx, input.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event details: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("event not found")
	}

	return &Output{result}, nil
}
