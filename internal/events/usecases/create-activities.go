package events

import (
	"context"
	"fmt"

	events "github.com/gabrielmatsan/checkin-gate/internal/events/domain"
	repository "github.com/gabrielmatsan/checkin-gate/internal/events/repository"
	identity "github.com/gabrielmatsan/checkin-gate/internal/identity/repository"
)

type CreateActivitiesUseCase struct {
	activityRepo repository.ActivityRepository
	eventRepo    repository.EventRepository
	userRepo     identity.UserRepository
}

type CreateActivitiesInput struct {
	UserID     string                     `json:"user_id"` // user ID that is creating the activities
	EventID    string                     `json:"event_id"`
	Activities []events.NewActivityParams `json:"activities"`
}

func NewCreateActivitiesUseCase(activityRepo repository.ActivityRepository, eventRepo repository.EventRepository, userRepo identity.UserRepository) *CreateActivitiesUseCase {
	return &CreateActivitiesUseCase{
		activityRepo: activityRepo,
		eventRepo:    eventRepo,
		userRepo:     userRepo,
	}
}

func (uc *CreateActivitiesUseCase) Execute(ctx context.Context, input CreateActivitiesInput) ([]*events.Activity, error) {

	// apenas 10 atividades por vez
	if len(input.Activities) > 10 {
		return nil, fmt.Errorf("only 10 activities at a time")
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.Role != "admin" {
		return nil, fmt.Errorf("user is not an admin")
	}

	// validar nomes duplicados no input
	nameSet := make(map[string]struct{}, len(input.Activities))
	names := make([]string, 0, len(input.Activities))

	for _, a := range input.Activities {
		// se o nome já existe, retornar erro
		if _, exists := nameSet[a.Name]; exists {
			return nil, fmt.Errorf("duplicate activity name in input: %s", a.Name)
		}

		// adicionar o nome ao set
		nameSet[a.Name] = struct{}{}
		// adicionar o nome ao slice
		names = append(names, a.Name)
	}

	// verifica se ja existem atividades com esses nomes no mesmo evento
	exisiting, err := uc.activityRepo.FindByEventIDAndNames(ctx, input.EventID, names)

	if err != nil {
		return nil, fmt.Errorf("failed to find activities by event ID and names: %w", err)
	}

	if len(exisiting) > 0 {
		return nil, fmt.Errorf("activities with the same names already exist for this event: %s", exisiting[0].Name)
	}

	// criar todas as atividades
	activities := make([]*events.Activity, 0, len(input.Activities))
	for _, a := range input.Activities {
		a.EventID = input.EventID
		activity, err := events.NewActivity(a)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	saved, err := uc.activityRepo.SaveAll(ctx, activities)
	if err != nil {
		return nil, err
	}

	return saved, nil
}
