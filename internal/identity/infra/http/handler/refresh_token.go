package handler

import (
	"encoding/json"
	"net/http"

	refreshtoken "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/refresh_token"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

// Request DTOs
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Response DTOs
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Handler
type RefreshTokenHandler struct {
	useCase *refreshtoken.UseCase
}

func NewRefreshTokenHandler(uc *refreshtoken.UseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{useCase: uc}
}

// Handle generates new access and refresh tokens.
// @Summary      Refresh tokens
// @Description  Uses a valid refresh token to obtain a new access token and rotated refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshTokenRequest   true  "Current refresh token"
// @Success      200   {object}  RefreshTokenResponse
// @Failure      400   {object}  lib.ErrorResponse
// @Failure      401   {object}  lib.ErrorResponse
// @Router       /auth/refresh [post]
func (h *RefreshTokenHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := lib.Validate(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := refreshTokenRequestToInput(&req, lib.GetClientIP(r), r.UserAgent())

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := refreshTokenOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// Mappers (internal to this handler)
func refreshTokenRequestToInput(req *RefreshTokenRequest, ipAddress, userAgent string) *refreshtoken.Input {
	return &refreshtoken.Input{
		RefreshToken: req.RefreshToken,
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
	}
}

func refreshTokenOutputToResponse(output *refreshtoken.Output) *RefreshTokenResponse {
	return &RefreshTokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}
}
