package http

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	checkinactivity "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/checkin_activity"
	createactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_activities"
	createevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_event"
	finishevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/finish_event"
	geteventdetails "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_details"
	geteventwithactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_with_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/http/handler"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/persistence"
	infraqueue "github.com/gabrielmatsan/checkin-gate/internal/events/infra/queue"
	eventsvc "github.com/gabrielmatsan/checkin-gate/internal/events/infra/service"
	identitypersistence "github.com/gabrielmatsan/checkin-gate/internal/identity/infra/persistence"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/service"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func RegisterEventsRoutes(r chi.Router, db *sqlx.DB, redisClient *redis.Client, cfg *config.Config, logger *zap.Logger) {
	jwtService := service.NewJWTService(cfg.JWTSecret)

	eventRepo := persistence.NewPostgresEventRepository(db)
	activityRepo := persistence.NewPostgresActivityRepository(db)
	checkInRepo := persistence.NewPostgresCheckInRepository(db)
	userRepo := identitypersistence.NewPostgresUserRepository(db)

	eventsTxProvider := persistence.NewPostgresTransactionProvider(db)

	userAuthSvc := eventsvc.NewUserAuthorizationAdapter(userRepo)
	certificateQueue := infraqueue.NewRedisCertificateQueue(redisClient)

	createEvent := createevent.NewUseCase(eventRepo, userAuthSvc)
	createActivities := createactivities.NewUseCase(activityRepo, eventRepo, userAuthSvc)
	getEventWithActivities := geteventwithactivities.NewUseCase(eventRepo, activityRepo)
	getEventDetails := geteventdetails.NewUseCase(eventRepo)
	checkInActivity := checkinactivity.NewUseCase(checkInRepo, activityRepo, eventRepo, userAuthSvc)
	finishEvent := finishevent.NewUseCase(eventsTxProvider, eventRepo, activityRepo, checkInRepo, userAuthSvc, certificateQueue)

	// Create individual handlers
	createEventHandler := handler.NewCreateEventHandler(logger, createEvent)
	createActivitiesHandler := handler.NewCreateActivitiesHandler(logger, createActivities)
	getEventWithActivitiesHandler := handler.NewGetEventWithActivitiesHandler(getEventWithActivities)
	getEventDetailsHandler := handler.NewGetEventDetailsHandler(getEventDetails)
	checkInActivityHandler := handler.NewCheckInActivityHandler(checkInActivity)
	finishEventHandler := handler.NewFinishEventHandler(finishEvent)

	r.Route("/events", func(r chi.Router) {
		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))

			r.Post("/", createEventHandler.Handle)
			r.Post("/activities", createActivitiesHandler.Handle)
			r.Get("/{event_id}/activities", getEventWithActivitiesHandler.Handle)
			r.Get("/{event_id}/details", getEventDetailsHandler.Handle)
			r.Post("/{event_id}/finish", finishEventHandler.Handle)
		})
	})

	r.Route("/activities", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))

			r.Post("/{activity_id}/checkin", checkInActivityHandler.Handle)
		})
	})
}
