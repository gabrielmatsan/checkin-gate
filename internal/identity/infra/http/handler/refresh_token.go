package handler

import (
	"net/http"

	refreshtoken "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/refresh_token"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

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
// @Produce      json
// @Success      200   {object}  RefreshTokenResponse
// @Failure      401   {object}  lib.ErrorResponse
// @Router       /auth/refresh [post]
func (h *RefreshTokenHandler) Handle(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, "missing refresh_token cookie")
		return
	}

	input := refreshTokenCookieToInput(cookie.Value, lib.GetClientIP(r), r.UserAgent())

	output, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set new Access Token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    output.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   900, // 15 minutes
	})

	// Set new Refresh Token cookie (rotation)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    output.RefreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   604800, // 7 days
	})

	resp := refreshTokenOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// Mappers (internal to this handler)
func refreshTokenCookieToInput(refreshToken, ipAddress, userAgent string) *refreshtoken.Input {
	return &refreshtoken.Input{
		RefreshToken: refreshToken,
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
