package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielmatsan/checkin-gate/internal/identity/handler/dto"
	usecases "github.com/gabrielmatsan/checkin-gate/internal/identity/usecases"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)



type AuthHandler struct {
	authenticateWithGoogle *usecases.AuthenticateWithGoogleUseCase
	refreshToken           *usecases.RefreshTokenUseCase
}

func NewAuthHandler(
	authenticateWithGoogle *usecases.AuthenticateWithGoogleUseCase,
	refreshToken *usecases.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		authenticateWithGoogle: authenticateWithGoogle,
		refreshToken:           refreshToken,
	}
}


// GoogleCallback authenticates a user via Google OAuth.
// @Summary      Authenticate with Google
// @Description  Exchanges a Google authorization code for access and refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GoogleCallbackRequest   true  "Google authorization code"
// @Success      200   {object}  dto.GoogleCallbackResponse
// @Failure      400   {object}  lib.ErrorResponse
// @Failure      401   {object}  lib.ErrorResponse
// @Router       /auth/google/callback [post]
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	var req dto.GoogleCallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !lib.ValidateAndRespond(w, &req) {
		return
	}

	// Request DTO → Use Case Input
	input := &usecases.AuthenticateWithGoogleInput{
		Code:      req.Code,
		IpAddress: lib.GetClientIP(r),
		UserAgent: r.UserAgent(),
	}

	output, err := h.authenticateWithGoogle.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Use Case Output → Response DTO
	resp := dto.GoogleCallbackResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		User: dto.UserResponse{
			ID:    output.User.ID,
			Email: output.User.Email,
			Name:  output.User.FirstName + " " + output.User.LastName,
			Role:  string(output.User.Role),
		},
	}

	lib.RespondJSON(w, http.StatusOK, resp)
}

// Refresh generates new access and refresh tokens.
// @Summary      Refresh tokens
// @Description  Uses a valid refresh token to obtain a new access token and rotated refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RefreshTokenRequest   true  "Current refresh token"
// @Success      200   {object}  dto.RefreshTokenResponse
// @Failure      400   {object}  lib.ErrorResponse
// @Failure      401   {object}  lib.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !lib.ValidateAndRespond(w, &req) {
		return
	}

	// Request DTO → Use Case Input
	input := usecases.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
		IpAddress:    lib.GetClientIP(r),
		UserAgent:    r.UserAgent(),
	}

	output, err := h.refreshToken.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Use Case Output → Response DTO
	resp := dto.RefreshTokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}

	lib.RespondJSON(w, http.StatusOK, resp)
}
