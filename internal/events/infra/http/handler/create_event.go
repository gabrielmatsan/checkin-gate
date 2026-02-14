package handler

import (
	"encoding/json"
	"net/http"
	"time"

	createevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_event"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"go.uber.org/zap"
)

// Request DTOs
type CreateEventRequest struct {
	Name           string    `json:"name" validate:"required,min=3,max=255"`
	AllowedDomains []string  `json:"allowed_domains" validate:"dive,fqdn"`
	Description    *string   `json:"description" validate:"omitempty,max=500"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

// Response DTOs
type CreateEventResponse struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	AllowedDomains []string   `json:"allowed_domains"`
	Description    *string    `json:"description,omitempty"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// Handler
type CreateEventHandler struct {
	logger  *zap.Logger
	useCase *createevent.UseCase
}

func NewCreateEventHandler(logger *zap.Logger, uc *createevent.UseCase) *CreateEventHandler {
	return &CreateEventHandler{
		logger:  logger,
		useCase: uc,
	}
}

// Handle creates a new event.
// @Summary      Create event
// @Description  Creates a new event. Only admins can create events.
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        body  body      CreateEventRequest  true  "Event details"
// @Success      201   {object}  CreateEventResponse
// @Failure      400   {object}  lib.ErrorResponse  "Invalid request body or validation error"
// @Failure      401   {object}  lib.ErrorResponse  "Unauthorized"
// @Failure      403   {object}  lib.ErrorResponse  "User not authorized to create event"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /events [post]
func (h *CreateEventHandler) Handle(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())


	var req CreateEventRequest
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

	input := createEventRequestToInput(&req, role)
	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		if err.Error() == "user not authorized to create event" {
			lib.RespondError(w, http.StatusForbidden, err.Error())
			return
		}
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := eventToCreateEventResponse(output.Event)
	lib.RespondJSON(w, http.StatusCreated, resp)
}

// Mappers (internal to this handler)
func createEventRequestToInput(req *CreateEventRequest, userRole string) *createevent.Input {
	return &createevent.Input{
		UserRole:       userRole,
		Name:           req.Name,
		AllowedDomains: &req.AllowedDomains,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}
}

func eventToCreateEventResponse(event *entity.Event) CreateEventResponse {
	return CreateEventResponse{
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
