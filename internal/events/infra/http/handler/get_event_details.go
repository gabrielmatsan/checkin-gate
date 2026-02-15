package handler

import (
	"net/http"
	"time"

	geteventdetails "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/get_event_details"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
)

// Response DTOs
type EventDetailsResponse struct {
	ID             string                         `json:"id"`
	Name           string                         `json:"name"`
	AllowedDomains []string                       `json:"allowed_domains"`
	Description    *string                        `json:"description,omitempty"`
	StartDate      time.Time                      `json:"start_date"`
	EndDate        time.Time                      `json:"end_date"`
	CreatedAt      time.Time                      `json:"created_at"`
	UpdatedAt      *time.Time                     `json:"updated_at,omitempty"`
	Activities     []ActivityWithCheckInsResponse `json:"activities"`
}

type ActivityWithCheckInsResponse struct {
	ActivityID   string            `json:"activity_id"`
	ActivityName string            `json:"activity_name"`
	CheckIns     []CheckInResponse `json:"check_ins"`
}

type CheckInResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ActivityID string    `json:"activity_id"`
	CheckedAt  time.Time `json:"checked_at"`
}

// Handler
type GetEventDetailsHandler struct {
	useCase *geteventdetails.UseCase
}

func NewGetEventDetailsHandler(uc *geteventdetails.UseCase) *GetEventDetailsHandler {
	return &GetEventDetailsHandler{useCase: uc}
}

// Handle gets event with activities and check-ins.
// @Summary      Get event details
// @Description  Gets an event with all activities and their check-ins. Admin only.
// @Tags         Events
// @Produce      json
// @Param        event_id  path      string  true  "Event ID"
// @Success      200   {object}  EventDetailsResponse
// @Failure      403   {object}  lib.ErrorResponse  "User not authorized"
// @Failure      404   {object}  lib.ErrorResponse  "Event not found"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Router       /events/{event_id}/details [get]
func (h *GetEventDetailsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	role := middleware.GetRole(r.Context())

	input := &geteventdetails.Input{
		EventID: eventID,
		Role:    role,
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		if err.Error() == "event not found" {
			lib.RespondError(w, http.StatusNotFound, err.Error())
			return
		}
		if err.Error() == "user not authorized" {
			lib.RespondError(w, http.StatusForbidden, err.Error())
			return
		}
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := toEventDetailsResponse(output.EventWithActivitiesAndCheckIns)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// Mappers
func toEventDetailsResponse(data *repository.EventWithActivitiesAndCheckIns) *EventDetailsResponse {
	activities := make([]ActivityWithCheckInsResponse, len(data.Activities))
	for i, a := range data.Activities {
		checkIns := make([]CheckInResponse, len(a.CheckIns))
		for j, c := range a.CheckIns {
			checkIns[j] = CheckInResponse{
				ID:         c.ID,
				UserID:     c.UserID,
				ActivityID: c.ActivityID,
				CheckedAt:  c.CheckedAt,
			}
		}
		activities[i] = ActivityWithCheckInsResponse{
			ActivityID:   a.ActivityID,
			ActivityName: a.ActivityName,
			CheckIns:     checkIns,
		}
	}

	return &EventDetailsResponse{
		ID:             data.Event.ID,
		Name:           data.Event.Name,
		AllowedDomains: data.Event.AllowedDomains,
		Description:    data.Event.Description,
		StartDate:      data.Event.StartDate,
		EndDate:        data.Event.EndDate,
		CreatedAt:      data.Event.CreatedAt,
		UpdatedAt:      data.Event.UpdatedAt,
		Activities:     activities,
	}
}
