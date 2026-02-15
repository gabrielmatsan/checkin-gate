package handler

import (
	"net/http"
	"time"

	checkinactivity "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/checkin_activity"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
)

// Response DTOs
type CheckInActivityResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ActivityID string    `json:"activity_id"`
	CheckedAt  time.Time `json:"checked_at"`
}

// Handler
type CheckInActivityHandler struct {
	useCase *checkinactivity.UseCase
}

func NewCheckInActivityHandler(uc *checkinactivity.UseCase) *CheckInActivityHandler {
	return &CheckInActivityHandler{useCase: uc}
}

// Handle performs a check-in to an activity.
// @Summary      Check-in to activity
// @Description  Performs a check-in to an activity. User must be authenticated.
// @Tags         CheckIn
// @Produce      json
// @Param        activity_id  path      string  true  "Activity ID"
// @Success      201   {object}  CheckInActivityResponse
// @Failure      400   {object}  lib.ErrorResponse  "Already checked in or outside activity time"
// @Failure      404   {object}  lib.ErrorResponse  "Activity not found"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Router       /activities/{activity_id}/checkin [post]
func (h *CheckInActivityHandler) Handle(w http.ResponseWriter, r *http.Request) {
	activityID := chi.URLParam(r, "activity_id")
	userID := middleware.GetUserID(r.Context())

	input := &checkinactivity.Input{
		UserID:     userID,
		ActivityID: activityID,
	}

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		switch err.Error() {
		case "activity not found":
			lib.RespondError(w, http.StatusNotFound, err.Error())
		case "user already checked in", "check-in not allowed outside activity time":
			lib.RespondError(w, http.StatusBadRequest, err.Error())
		default:
			lib.RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := &CheckInActivityResponse{
		ID:         output.CheckIn.ID,
		UserID:     output.CheckIn.UserID,
		ActivityID: output.CheckIn.ActivityID,
		CheckedAt:  output.CheckIn.CheckedAt,
	}

	lib.RespondJSON(w, http.StatusCreated, resp)
}
