package handler

import (
	"encoding/json"
	"net/http"
	"time"

	createactivities "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/create_activities"
	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"go.uber.org/zap"
)

type CreateActivityItem struct {
	Name        string    `json:"name" validate:"required"`
	Description *string   `json:"description"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type CreateActivitiesRequest struct {
	EventID    string               `json:"event_id" validate:"required"`
	Activities []CreateActivityItem `json:"activities" validate:"required,dive"`
}

type CreateActivityResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	EventID     string     `json:"event_id"`
	Description *string    `json:"description,omitempty"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type CreateActivitiesHandler struct {
	useCase *createactivities.UseCase
	logger  *zap.Logger
}

func NewCreateActivitiesHandler(logger *zap.Logger, uc *createactivities.UseCase) *CreateActivitiesHandler {
	return &CreateActivitiesHandler{logger: logger, useCase: uc}
}

// Handle creates activities for an event.
// @Summary      Create activities
// @Description  Creates one or more activities for an event. User must be authenticated.
// @Tags         Activities
// @Accept       json
// @Produce      json
// @Param        request  body      CreateActivitiesRequest  true  "Activities to create"
// @Success      201      {array}   CreateActivityResponse
// @Failure      400      {object}  lib.ErrorResponse  "Invalid request body or validation error"
// @Failure      500      {object}  lib.ErrorResponse  "Internal server error"
// @Router       /events/activities [post]
func (h *CreateActivitiesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req CreateActivitiesRequest

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

	input := createActivitiesRequestToInput(&req, userID)

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := activitiesToCreateActivitiesResponse(output.Activities)
	lib.RespondJSON(w, http.StatusCreated, resp)
}

// Mappers

// a partir de CreateActivitiesRequest, cria o inputs Activities
func activityItemsToActivityInputs(items []CreateActivityItem) []createactivities.ActivityInput {
	inputs := make([]createactivities.ActivityInput, len(items))
	for i, item := range items {
		inputs[i] = createactivities.ActivityInput{
			Name:        item.Name,
			Description: item.Description,
			StartDate:   item.StartDate,
			EndDate:     item.EndDate,
		}
	}
	return inputs
}

func createActivitiesRequestToInput(req *CreateActivitiesRequest, userID string) *createactivities.Input {
	return &createactivities.Input{
		UserID:     userID,
		EventID:    req.EventID,
		Activities: activityItemsToActivityInputs(req.Activities),
	}
}

func activitiesToCreateActivitiesResponse(activities []*entity.Activity) []CreateActivityResponse {
	responses := make([]CreateActivityResponse, len(activities))
	for i, activity := range activities {
		responses[i] = CreateActivityResponse{
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
	return responses
}
