package events

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielmatsan/checkin-gate/internal/events/handler/dto"
	events "github.com/gabrielmatsan/checkin-gate/internal/events/usecases"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
)

type EventHandler struct {
	createEvent *events.CreateEventUseCase
}

func NewEventHandler(createEvent *events.CreateEventUseCase) *EventHandler {
	return &EventHandler{
		createEvent: createEvent,
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())

	isUserAuthorized := middleware.ValidateRole(role, []string{"admin"})

	if !isUserAuthorized {
		lib.RespondError(w, http.StatusForbidden, "user not authorized to create event")
		return
	}

	var req dto.CreateEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := lib.Validate(&req); err != nil {
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
