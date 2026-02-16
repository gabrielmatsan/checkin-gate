package handler

import (
	"log"
	"net/http"

	finishevent "github.com/gabrielmatsan/checkin-gate/internal/events/application/usecase/finish_event"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
)

type FinishEventHandler struct {
	useCase *finishevent.UseCase
}

func NewFinishEventHandler(uc *finishevent.UseCase) *FinishEventHandler {
	return &FinishEventHandler{useCase: uc}
}

// Handle finishes an event and enqueues certificate jobs.
// @Summary      Finish event
// @Description  Finishes an event and enqueues certificate generation jobs for all check-ins. Admin only.
// @Tags         Events
// @Produce      json
// @Param        event_id  path      string  true  "Event ID"
// @Success      202   {object}  map[string]string  "Jobs enqueued"
// @Failure      400   {object}  lib.ErrorResponse  "Activities not ended or no check-ins"
// @Failure      403   {object}  lib.ErrorResponse  "User is not an admin"
// @Failure      404   {object}  lib.ErrorResponse  "Event not found"
// @Failure      500   {object}  lib.ErrorResponse  "Internal server error"
// @Router       /events/{event_id}/finish [post]
func (h *FinishEventHandler) Handle(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	userID := middleware.GetUserID(r.Context())

	log.Printf("[finish_event] eventID=%q userID=%q", eventID, userID)

	input := &finishevent.Input{
		EventID: eventID,
		UserID:  userID,
	}

	err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		log.Printf("[finish_event] error=%q", err.Error())
		switch err.Error() {
		case "user not found", "event not found":
			lib.RespondError(w, http.StatusNotFound, err.Error())
		case "user is not an admin":
			lib.RespondError(w, http.StatusForbidden, err.Error())
		case "activity has not ended", "no activities found for event", "no check-ins found for activities":
			lib.RespondError(w, http.StatusBadRequest, err.Error())
		default:
			lib.RespondError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	lib.RespondJSON(w, http.StatusAccepted, map[string]string{
		"message": "certificate jobs enqueued",
	})
}
