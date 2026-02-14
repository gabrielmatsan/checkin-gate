package checkinactivity

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
)

type Input struct {
	UserID     string
	ActivityID string
}

type Output struct {
	CheckIn *entity.CheckIn
}

type UseCase struct {
	checkInRepo  repository.CheckInRepository
	activityRepo repository.ActivityRepository
}

func NewUseCase(
	checkInRepo repository.CheckInRepository,
	activityRepo repository.ActivityRepository,
) *UseCase {
	return &UseCase{
		checkInRepo:  checkInRepo,
		activityRepo: activityRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	// 1. Verificar se atividade existe
	activity, err := uc.activityRepo.FindByID(ctx, input.ActivityID)
	if err != nil {
		return nil, fmt.Errorf("failed to find activity: %w", err)
	}
	if activity == nil {
		return nil, fmt.Errorf("activity not found")
	}

	// 2. Verificar se já fez check-in
	existing, err := uc.checkInRepo.FindByUserAndActivity(ctx, input.UserID, input.ActivityID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing check-in: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("user already checked in")
	}

	// 3. Verificar se está no horário da atividade
	now := time.Now()
	if now.Before(activity.StartDate) || now.After(activity.EndDate) {
		return nil, fmt.Errorf("check-in not allowed outside activity time")
	}

	// 4. Criar check-in
	checkIn, err := entity.NewCheckIn(entity.NewCheckInParams{
		UserID:     input.UserID,
		ActivityID: input.ActivityID,
	})
	if err != nil {
		return nil, err
	}

	// 5. Salvar
	saved, err := uc.checkInRepo.Save(ctx, checkIn)
	if err != nil {
		return nil, fmt.Errorf("failed to save check-in: %w", err)
	}

	return &Output{CheckIn: saved}, nil
}
