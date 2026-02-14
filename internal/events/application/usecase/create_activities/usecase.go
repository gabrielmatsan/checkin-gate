package createactivities

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/service"
)

type ActivityInput struct {
	Name        string
	Description *string
	StartDate   time.Time
	EndDate     time.Time
}

type Input struct {
	UserID     string
	EventID    string
	Activities []ActivityInput
}

type Output struct {
	Activities []*entity.Activity
}

type UseCase struct {
	activityRepo repository.ActivityRepository
	eventRepo    repository.EventRepository
	userAuthSvc  service.UserAuthorizationService
}

func NewUseCase(
	activityRepo repository.ActivityRepository,
	eventRepo repository.EventRepository,
	userAuthSvc service.UserAuthorizationService,
) *UseCase {
	return &UseCase{
		activityRepo: activityRepo,
		eventRepo:    eventRepo,
		userAuthSvc:  userAuthSvc,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	// Only 10 activities at a time
	if len(input.Activities) > 10 {
		return nil, fmt.Errorf("only 10 activities at a time")
	}

	// Check if user is admin using the authorization service
	isAdmin, err := uc.userAuthSvc.IsUserAdmin(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, fmt.Errorf("user is not an admin")
	}

	// Validate for duplicate names in input
	nameSet := make(map[string]struct{}, len(input.Activities))
	names := make([]string, 0, len(input.Activities))

	for _, a := range input.Activities {
		if _, exists := nameSet[a.Name]; exists {
			return nil, fmt.Errorf("duplicate activity name in input: %s", a.Name)
		}
		nameSet[a.Name] = struct{}{}
		names = append(names, a.Name)
	}

	// Check if activities with these names already exist for this event
	existing, err := uc.activityRepo.FindByEventIDAndNames(ctx, input.EventID, names)
	if err != nil {
		return nil, fmt.Errorf("failed to find activities by event ID and names: %w", err)
	}
	if len(existing) > 0 {
		return nil, fmt.Errorf("activities with the same names already exist for this event: %s", existing[0].Name)
	}

	// Create all activities
	activities := make([]*entity.Activity, 0, len(input.Activities))
	for _, a := range input.Activities {
		activity, err := entity.NewActivity(entity.NewActivityParams{
			Name:        a.Name,
			EventID:     input.EventID,
			Description: a.Description,
			StartDate:   a.StartDate,
			EndDate:     a.EndDate,
		})
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	saved, err := uc.activityRepo.SaveAll(ctx, activities)
	if err != nil {
		return nil, err
	}

	return &Output{
		Activities: saved,
	}, nil
}
