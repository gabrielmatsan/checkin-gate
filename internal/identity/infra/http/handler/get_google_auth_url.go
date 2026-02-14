package handler

import (
	"net/http"

	getgoogleauthurl "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/get_google_auth_url"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

// Response DTOs
type GetGoogleAuthURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"`
}

// Handler
type GetGoogleAuthURLHandler struct {
	useCase *getgoogleauthurl.UseCase
}

func NewGetGoogleAuthURLHandler(uc *getgoogleauthurl.UseCase) *GetGoogleAuthURLHandler {
	return &GetGoogleAuthURLHandler{useCase: uc}
}

// Handle generates a new Google OAuth URL.
// @Summary      Get Google OAuth URL
// @Description  Generates a new Google OAuth URL for authentication
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  GetGoogleAuthURLResponse
// @Failure      500   {object}  lib.ErrorResponse
// @Router       /auth/google/url [get]
func (h *GetGoogleAuthURLHandler) Handle(w http.ResponseWriter, r *http.Request) {
	output, err := h.useCase.Execute()
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := googleAuthURLOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// Mappers (internal to this handler)
func googleAuthURLOutputToResponse(output *getgoogleauthurl.Output) *GetGoogleAuthURLResponse {
	return &GetGoogleAuthURLResponse{
		URL:   output.URL,
		State: output.State,
	}
}
