package handler

import (
	"encoding/json"
	"net/http"

	createevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_event"
	geteventwithactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_with_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/http/dto"
	"github.com/gabrielmatsan/checkin-gate/internal/events/infra/http/mapper"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type EventHandler struct {
	logger                 *zap.Logger
	createEvent            *createevent.UseCase
	getEventWithActivities *geteventwithactivities.UseCase
}

func NewEventHandler(
	logger *zap.Logger,
	createEvent *createevent.UseCase,
	getEventWithActivities *geteventwithactivities.UseCase,
) *EventHandler {
	return &EventHandler{
		logger:                 logger,
		createEvent:            createEvent,
		getEventWithActivities: getEventWithActivities,
	}
}

// CreateEvent creates a new event.
// @Summary      Create event
// @Description  Creates a new event. Only admins can create events.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateEventRequest  true  "Event details"
// @Success      201   {object}  dto.EventResponse
// @Failure      400   {object}  lib.ErrorResponse  "Invalid request body or validation error"
// @Failure      401   {object}  lib.ErrorResponse  "Unauthorized"
// @Failure      403   {object}  lib.ErrorResponse  "User not authorized to create event"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /events [post]
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	isUserAuthorized := middleware.ValidateRole(role, []string{"admin"})
	if !isUserAuthorized {
		lib.RespondError(w, http.StatusForbidden, "user not authorized to create event")
		return
	}

	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body",
			zap.Error(err),
			zap.String("content_type", r.Header.Get("Content-Type")),
		)
		lib.RespondError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	h.logger.Info("request decoded", zap.Any("request", req))

	if err := lib.Validate(&req); err != nil {
		h.logger.Error("validation failed", zap.Error(err))
		lib.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := mapper.CreateEventRequestToInput(&req)
	output, err := h.createEvent.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := mapper.EventToResponse(output.Event)
	lib.RespondJSON(w, http.StatusCreated, resp)
}

// GetEventWithActivities gets an event with its activities.
// @Summary      Get event with activities
// @Description  Gets an event with its activities.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event_id  path      string  true  "Event ID"
// @Success      200   {object}  dto.EventWithActivitiesResponse
// @Failure      400   {object}  lib.ErrorResponse  "Invalid request body or validation error"
// @Failure      404   {object}  lib.ErrorResponse  "Event not found"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /events/{event_id}/activities [get]
func (h *EventHandler) GetEventWithActivities(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")

	input := &geteventwithactivities.Input{
		EventID: eventID,
	}

	output, err := h.getEventWithActivities.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := mapper.GetEventWithActivitiesOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}
