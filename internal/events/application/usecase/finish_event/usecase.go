package finishevent

import (
	"context"
	"errors"
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/queue"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/service"
	"golang.org/x/sync/errgroup"
)

type Input struct {
	EventID string `json:"event_id" validate:"required"`
	UserID  string `json:"user_id" validate:"required"`
}

type UseCase struct {
	eventRepo        repository.EventRepository
	activityRepo     repository.ActivityRepository
	checkInRepo      repository.CheckInRepository
	userAuthSvc      service.UserAuthorizationService
	certificateQueue queue.CertificateQueue
}

func NewUseCase(eventRepo repository.EventRepository, activityRepo repository.ActivityRepository, checkInRepo repository.CheckInRepository, userAuthSvc service.UserAuthorizationService, certificateQueue queue.CertificateQueue) *UseCase {
	return &UseCase{
		eventRepo:        eventRepo,
		activityRepo:     activityRepo,
		checkInRepo:      checkInRepo,
		userAuthSvc:      userAuthSvc,
		certificateQueue: certificateQueue,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *Input) error {

	user, err := uc.userAuthSvc.GetUserByID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %w", err)
	}

	if user == nil {
		return errors.New("user not found")
	}

	if !user.IsAdmin {
		return errors.New("user is not an admin")
	}

	var (
		event      *entity.Event
		activities []*entity.Activity
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		event, err = uc.eventRepo.FindByID(gCtx, input.EventID)
		if err != nil {
			return fmt.Errorf("failed to find event by ID: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		activities, err = uc.activityRepo.FindByEventID(gCtx, input.EventID)
		if err != nil {
			return fmt.Errorf("failed to find activities by event ID: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if event == nil {
		return errors.New("event not found")
	}

	if len(activities) == 0 {
		return errors.New("no activities found for event")
	}

	activityIDs := make([]string, len(activities))
	for i, activity := range activities {
		activityIDs[i] = activity.ID
		if !activity.HasEnded() {
			return errors.New("activity has not ended")
		}
	}

	// buscar check-ins das atividades
	checkIns, err := uc.checkInRepo.FindByActivityIDs(ctx, activityIDs)

	if err != nil {
		return fmt.Errorf("failed to find check-ins by activity IDs: %w", err)
	}

	if len(checkIns) == 0 {
		return errors.New("no check-ins found for activities")
	}

	// faz array com userIDs de checkIns, removendo duplicados
	seen := make(map[string]struct{}, len(checkIns))
	userIDs := make([]string, 0, len(checkIns))

	for _, checkIn := range checkIns {
		// se o userID ja foi visto, nao adiciona novamente
		if _, ok := seen[checkIn.UserID]; !ok {
			seen[checkIn.UserID] = struct{}{}
			userIDs = append(userIDs, checkIn.UserID)
		}
	}

	users, err := uc.userAuthSvc.GetUserInfoBatch(ctx, userIDs)
	if err != nil {
		return fmt.Errorf("failed to get user info batch: %w", err)
	}

	if len(users) == 0 {
		return errors.New("no users found for check-ins")
	}

	// indexar users por userID
	userIndex := make(map[string]*service.UserInfo)
	for _, user := range users {
		userIndex[user.ID] = user
	}

	// indexar activities por activityID
	activityIndex := make(map[string]*entity.Activity)
	for _, activity := range activities {
		activityIndex[activity.ID] = activity
	}

	// agrupar checkIns com informacoes de users
	jobs := make([]*queue.CertificateJob, 0, len(checkIns))
	for _, checkIn := range checkIns {

		user, ok := userIndex[checkIn.UserID]
		if !ok {
			return fmt.Errorf("user not found: %s", checkIn.UserID)
		}

		activity, ok := activityIndex[checkIn.ActivityID]
		if !ok {
			return fmt.Errorf("activity not found: %s", checkIn.ActivityID)
		}

		job, err := queue.NewCertificateJob(queue.NewCertificateJobParams{
			EventInfo: queue.EventInfo{
				EventID:   event.ID,
				EventName: event.Name,
			},
			UserInfo: queue.UserInfo{
				UserID:    user.ID,
				UserName:  user.FirstName + " " + user.LastName,
				UserEmail: user.Email,
			},
			ActivityInfo: queue.ActivityInfo{
				ActivityID:   activity.ID,
				ActivityName: activity.Name,
				ActivityDate: activity.StartDate,
				StartTime:    activity.StartDate,
				EndTime:      activity.EndDate,
			},
			CheckedAt: checkIn.CheckedAt,
		})
		if err != nil {
			return fmt.Errorf("failed to create certificate job: %w", err)
		}

		jobs = append(jobs, job)
	}

	// enfileirar jobs
	if err := uc.certificateQueue.EnqueueBatch(ctx, jobs); err != nil {
		return fmt.Errorf("failed to enqueue certificate jobs: %w", err)
	}

	return nil
}
