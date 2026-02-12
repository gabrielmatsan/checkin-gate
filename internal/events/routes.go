package events

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	handler "github.com/gabrielmatsan/checkin-gate/internal/events/handler"
	repository "github.com/gabrielmatsan/checkin-gate/internal/events/infra"
	usecases "github.com/gabrielmatsan/checkin-gate/internal/events/usecases"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	jwtService := infra.NewJWTService(cfg.JWTSecret)

	eventRepo := repository.NewPostgresEventRepository(db)
	activityRepo := repository.NewPostgresActivityRepository(db)

	createEvent := usecases.NewCreateEventUseCase(eventRepo)

	getEventWithActivities := usecases.NewGetEventWithActivitiesUseCase(eventRepo, activityRepo)

	eventHandler := handler.NewEventHandler(createEvent, getEventWithActivities)

	r.Route("/events", func(r chi.Router) {
		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))

			r.Post("/", eventHandler.CreateEvent)
			r.Get("/{event_id}/activities", eventHandler.GetEventWithActivities)
		})
	})
}
