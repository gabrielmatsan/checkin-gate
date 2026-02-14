package handler

import (
	"net/http"

	authenticatewithgoogle "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/authenticate_with_google"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/domain/entity"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

// Request DTOs
type GoogleCallbackRequest struct {
	Code  string `validate:"required"`
	State string `validate:"required"`
}

// Response DTOs
type GoogleCallbackResponse struct {
	AccessToken  string                     `json:"access_token"`
	RefreshToken string                     `json:"refresh_token"`
	User         GoogleCallbackUserResponse `json:"user"`
}

type GoogleCallbackUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

// Handler
type GoogleCallbackHandler struct {
	useCase *authenticatewithgoogle.UseCase
}

func NewGoogleCallbackHandler(uc *authenticatewithgoogle.UseCase) *GoogleCallbackHandler {
	return &GoogleCallbackHandler{useCase: uc}
}

// Handle authenticates a user via Google OAuth.
// @Summary      Authenticate with Google
// @Description  Exchanges a Google authorization code for access and refresh tokens
// @Tags         Auth
// @Produce      json
// @Param        code   query     string  true  "Google authorization code"
// @Param        state  query     string  true  "State parameter for CSRF protection"
// @Success      200   {object}  GoogleCallbackResponse
// @Failure      400   {object}  lib.ErrorResponse
// @Failure      401   {object}  lib.ErrorResponse
// @Router       /auth/google/callback [get]
func (h *GoogleCallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := GoogleCallbackRequest{
		Code:  r.URL.Query().Get("code"),
		State: r.URL.Query().Get("state"),
	}

	if err := lib.Validate(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := googleCallbackRequestToInput(&req, lib.GetClientIP(r), r.UserAgent())

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set Access Token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    output.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   900, // 15 minutes
	})

	// Set Refresh Token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})

	resp := authOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// Mappers (internal to this handler)
func googleCallbackRequestToInput(req *GoogleCallbackRequest, ipAddress, userAgent string) *authenticatewithgoogle.Input {
	return &authenticatewithgoogle.Input{
		Code:      req.Code,
		IpAddress: ipAddress,
		UserAgent: userAgent,
	}
}

func authOutputToResponse(output *authenticatewithgoogle.Output) *GoogleCallbackResponse {
	return &GoogleCallbackResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User:         userToResponse(output.User),
	}
}

func userToResponse(user *entity.User) GoogleCallbackUserResponse {
	return GoogleCallbackUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.FirstName + " " + user.LastName,
		Role:  string(user.Role),
	}
}
