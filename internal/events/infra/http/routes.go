package http

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	createevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_event"
	geteventwithactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_with_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/http/handler"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/persistence"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/service"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config, logger *zap.Logger) {
	jwtService := service.NewJWTService(cfg.JWTSecret)

	eventRepo := persistence.NewPostgresEventRepository(db)
	activityRepo := persistence.NewPostgresActivityRepository(db)

	createEvent := createevent.NewUseCase(eventRepo)
	getEventWithActivities := geteventwithactivities.NewUseCase(eventRepo, activityRepo)

	eventHandler := handler.NewEventHandler(logger, createEvent, getEventWithActivities)

	r.Route("/events", func(r chi.Router) {
		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))

			r.Post("/", eventHandler.CreateEvent)
			r.Get("/{event_id}/activities", eventHandler.GetEventWithActivities)
		})
	})
}
