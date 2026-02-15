package handler

import (
	"net/http"
	"time"

	geteventwithactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_with_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/go-chi/chi/v5"
)

// Response DTOs
type EventResponse struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	AllowedDomains []string   `json:"allowed_domains"`
	Description    *string    `json:"description,omitempty"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type ActivityResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	EventID     string     `json:"event_id"`
	Description *string    `json:"description,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type EventWithActivitiesResponse struct {
	Event      EventResponse      `json:"event"`
	Activities []ActivityResponse `json:"activities"`
}

// Handler
type GetEventWithActivitiesHandler struct {
	useCase *geteventwithactivities.UseCase
}

func NewGetEventWithActivitiesHandler(uc *geteventwithactivities.UseCase) *GetEventWithActivitiesHandler {
	return &GetEventWithActivitiesHandler{useCase: uc}
}

// Handle gets an event with its activities.
// @Summary      Get event with activities
// @Description  Gets an event with its activities.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event_id  path      string  true  "Event ID"
// @Success      200   {object}  EventWithActivitiesResponse
// @Failure      400   {object}  lib.ErrorResponse  "Invalid request body or validation error"
// @Failure      404   {object}  lib.ErrorResponse  "Event not found"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Router       /events/{event_id}/activities [get]
func (h *GetEventWithActivitiesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")

	input := &geteventwithactivities.Input{
		EventID: eventID,
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := getEventWithActivitiesOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// Mappers (internal to this handler)
func eventToResponse(event *entity.Event) EventResponse {
	return EventResponse{
		ID:             event.ID,
		Name:           event.Name,
		AllowedDomains: event.AllowedDomains,
		Description:    event.Description,
		StartDate:      event.StartDate,
		EndDate:        event.EndDate,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func activityToResponse(activity *entity.Activity) ActivityResponse {
	return ActivityResponse{
		ID:          activity.ID,
		Name:        activity.Name,
		EventID:     activity.EventID,
		Description: activity.Description,
		StartDate:   activity.StartDate,
		EndDate:     activity.EndDate,
		CreatedAt:   activity.CreatedAt,
		UpdatedAt:   activity.UpdatedAt,
	}
}

func activitiesToResponse(activities []*entity.Activity) []ActivityResponse {
	result := make([]ActivityResponse, len(activities))
	for i, a := range activities {
		result[i] = activityToResponse(a)
	}
	return result
}

func getEventWithActivitiesOutputToResponse(output *geteventwithactivities.Output) *EventWithActivitiesResponse {
	return &EventWithActivitiesResponse{
		Event:      eventToResponse(output.Event),
		Activities: activitiesToResponse(output.Activities),
	}
}
