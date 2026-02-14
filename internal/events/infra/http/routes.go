package http

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	checkinactivity "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/checkin_activity"
	createevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_event"
	geteventdetails "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_details"
	geteventwithactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_with_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/http/handler"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/persistence"
	eventsvc "github.com/gabrielmatsan/checkin-gate/internal/events/infra/service"
	identitypersistence "github.com/gabrielmatsan/checkin-gate/internal/identity/infra/persistence"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/service"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RegisterEventsRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config, logger *zap.Logger) {
	jwtService := service.NewJWTService(cfg.JWTSecret)

	eventRepo := persistence.NewPostgresEventRepository(db)
	activityRepo := persistence.NewPostgresActivityRepository(db)
	checkInRepo := persistence.NewPostgresCheckInRepository(db)
	userRepo := identitypersistence.NewPostgresUserRepository(db)

	userAuthSvc := eventsvc.NewUserAuthorizationAdapter(userRepo)

	createEvent := createevent.NewUseCase(eventRepo, userAuthSvc)
	getEventWithActivities := geteventwithactivities.NewUseCase(eventRepo, activityRepo)
	getEventDetails := geteventdetails.NewUseCase(eventRepo)
	checkInActivity := checkinactivity.NewUseCase(checkInRepo, activityRepo)

	// Create individual handlers
	createEventHandler := handler.NewCreateEventHandler(logger, createEvent)
	getEventWithActivitiesHandler := handler.NewGetEventWithActivitiesHandler(getEventWithActivities)
	getEventDetailsHandler := handler.NewGetEventDetailsHandler(getEventDetails)
	checkInActivityHandler := handler.NewCheckInActivityHandler(checkInActivity)

	r.Route("/events", func(r chi.Router) {
		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))

			r.Post("/", createEventHandler.Handle)
			r.Get("/{event_id}/activities", getEventWithActivitiesHandler.Handle)
			r.Get("/{event_id}/details", getEventDetailsHandler.Handle)
		})
	})

	r.Route("/activities", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))

			r.Post("/{activity_id}/checkin", checkInActivityHandler.Handle)
		})
	})
}
