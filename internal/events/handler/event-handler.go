package events

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielmatsan/checkin-gate/internal/events/handler/dto"
	events "github.com/gabrielmatsan/checkin-gate/internal/events/usecases"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	//"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type EventHandler struct {
	logger                 *zap.Logger
	createEvent            *events.CreateEventUseCase
	getEventWithActivities *events.GetEventWithActivitiesUseCase
}

func NewEventHandler(logger *zap.Logger, createEvent *events.CreateEventUseCase, getEventWithActivities *events.GetEventWithActivitiesUseCase) *EventHandler {
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
	//role := middleware.GetRole(r.Context())

	// isUserAuthorized := middleware.ValidateRole(role, []string{"admin"})

	// if !isUserAuthorized {
	// 	lib.RespondError(w, http.StatusForbidden, "user not authorized to create event")
	// 	return
	// }

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

	input := events.CreateEventInput{
		Name:           req.Name,
		AllowedDomains: &req.AllowedDomains,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}

	output, err := h.createEvent.Execute(r.Context(), input)

	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.EventResponse{
		ID:             output.ID,
		Name:           output.Name,
		AllowedDomains: output.AllowedDomains,
		Description:    output.Description,
		StartDate:      output.StartDate,
		EndDate:        output.EndDate,
		CreatedAt:      output.CreatedAt,
	}

	lib.RespondJSON(w, http.StatusCreated, resp)
}

// GetEventWithActivities gets an event with its activities.
// @Summary      Get event with activities
// @Description  Gets an event with its activities.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event_id  path      string  true  "Event ID"
// @Success      200   {object}  events.EventWithActivities
// @Failure      400   {object}  lib.ErrorResponse  "Invalid request body or validation error"
// @Failure      404   {object}  lib.ErrorResponse  "Event not found"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /events/{event_id}/activities [get]
func (h *EventHandler) GetEventWithActivities(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	var req dto.EventWithActivitiesRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := lib.Validate(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := events.GetEventWithActivitiesInput{
		EventID: eventID,
	}

	output, err := h.getEventWithActivities.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := events.EventWithActivities{
		Event:      output.Event,
		Activities: output.Activities,
	}

	lib.RespondJSON(w, http.StatusOK, resp)
}
