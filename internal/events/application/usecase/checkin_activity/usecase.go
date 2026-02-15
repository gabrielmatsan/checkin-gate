package checkinactivity

import (
	"context"
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/service"
	"golang.org/x/sync/errgroup"
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
	eventRepo    repository.EventRepository
	userAuthSvc  service.UserAuthorizationService
}

func NewUseCase(
	checkInRepo repository.CheckInRepository,
	activityRepo repository.ActivityRepository,
	eventRepo repository.EventRepository,
	userAuthSvc service.UserAuthorizationService,
) *UseCase {
	return &UseCase{
		checkInRepo:  checkInRepo,
		activityRepo: activityRepo,
		eventRepo:    eventRepo,
		userAuthSvc:  userAuthSvc,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) (*Output, error) {
	var activity *entity.Activity
	var existing *entity.CheckIn

	// 1. Buscar atividade e check-in existente em paralelo
	g1, gCtx := errgroup.WithContext(ctx)

	g1.Go(func() error {
		var err error
		activity, err = uc.activityRepo.FindByID(gCtx, input.ActivityID)
		if err != nil {
			return fmt.Errorf("failed to find activity: %w", err)
		}
		return nil
	})

	g1.Go(func() error {
		var err error
		existing, err = uc.checkInRepo.FindByUserAndActivity(gCtx, input.UserID, input.ActivityID)
		if err != nil {
			return fmt.Errorf("failed to check existing check-in: %w", err)
		}
		return nil
	})

	if err := g1.Wait(); err != nil {
		return nil, err
	}

	// 2. Validar atividade existe
	if activity == nil {
		return nil, fmt.Errorf("activity not found")
	}

	// 3. Verificar se já fez check-in
	if existing != nil {
		return nil, fmt.Errorf("user already checked in")
	}

	// 4. Verificar se está no horário da atividade
	// TODO: descomentar após testes
	// if !activity.IsCheckInAllowed(time.Now()) {
	// 	return nil, fmt.Errorf("check-in not allowed outside activity time")
	// }

	// 5. Busca evento e email do usuário em paralelo
	var event *entity.Event
	var userEmail string

	g2, gCtx2 := errgroup.WithContext(ctx)

	g2.Go(func() error {
		var err error
		event, err = uc.eventRepo.FindByID(gCtx2, activity.EventID)
		if err != nil {
			return fmt.Errorf("failed to find event: %w", err)
		}
		return nil
	})

	g2.Go(func() error {
		var err error
		userEmail, err = uc.userAuthSvc.GetUserEmail(gCtx2, input.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user email: %w", err)
		}
		return nil
	})

	if err := g2.Wait(); err != nil {
		return nil, err
	}

	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	if userEmail == "" {
		return nil, fmt.Errorf("user email not found")
	}

	// 6. Verificar se o dominio do usuario está permitido
	if !event.IsAllowedDomain(userEmail) {
		return nil, fmt.Errorf("user domain not allowed")
	}

	// 7. Criar check-in
	checkIn, err := entity.NewCheckIn(entity.NewCheckInParams{
		UserID:     input.UserID,
		ActivityID: input.ActivityID,
	})
	if err != nil {
		return nil, err
	}

	// 8. Salvar
	saved, err := uc.checkInRepo.Save(ctx, checkIn)
	if err != nil {
		return nil, fmt.Errorf("failed to save check-in: %w", err)
	}

	return &Output{CheckIn: saved}, nil
}
