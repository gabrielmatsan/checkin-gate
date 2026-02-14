package handler

import (
	"encoding/json"
	"net/http"

	authenticatewithgoogle "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/authenticate_with_google"
	getgoogleauthurl "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/get_google_auth_url"
	refreshtoken "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/refresh_token"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/http/dto"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/http/mapper"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
)

type AuthHandler struct {
	authenticateWithGoogle *authenticatewithgoogle.UseCase
	refreshToken           *refreshtoken.UseCase
	getGoogleAuthURL       *getgoogleauthurl.UseCase
}

func NewAuthHandler(
	getGoogleAuthURL *getgoogleauthurl.UseCase,
	authenticateWithGoogle *authenticatewithgoogle.UseCase,
	refreshToken *refreshtoken.UseCase,
) *AuthHandler {
	return &AuthHandler{
		getGoogleAuthURL:       getGoogleAuthURL,
		authenticateWithGoogle: authenticateWithGoogle,
		refreshToken:           refreshToken,
	}
}

// GoogleCallback authenticates a user via Google OAuth.
// @Summary      Authenticate with Google
// @Description  Exchanges a Google authorization code for access and refresh tokens
// @Tags         Auth
// @Produce      json
// @Param        code   query     string  true  "Google authorization code"
// @Param        state  query     string  true  "State parameter for CSRF protection"
// @Success      200   {object}  dto.GoogleCallbackResponse
// @Failure      400   {object}  lib.ErrorResponse
// @Failure      401   {object}  lib.ErrorResponse
// @Router       /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	req := dto.GoogleCallbackRequest{
		Code:  r.URL.Query().Get("code"),
		State: r.URL.Query().Get("state"),
	}

	if err := lib.Validate(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := mapper.GoogleCallbackRequestToInput(&req, lib.GetClientIP(r), r.UserAgent())

	output, err := h.authenticateWithGoogle.Execute(r.Context(), input)
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

	resp := mapper.AuthOutputToResponse(output)
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

	if err := lib.Validate(&req); err != nil {
		lib.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := mapper.RefreshTokenRequestToInput(&req, lib.GetClientIP(r), r.UserAgent())

	output, err := h.refreshToken.Execute(r.Context(), input)
	if err != nil {
		lib.RespondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := mapper.RefreshTokenOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}

// GetGoogleAuthURL generates a new Google OAuth URL.
// @Summary      Get Google OAuth URL
// @Description  Generates a new Google OAuth URL for authentication
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  dto.GoogleAuthURLResponse
// @Failure      500   {object}  lib.ErrorResponse
// @Router       /auth/google/url [get]
func (h *AuthHandler) GetGoogleAuthURL(w http.ResponseWriter, r *http.Request) {
	output, err := h.getGoogleAuthURL.Execute()
	if err != nil {
		lib.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := mapper.GoogleAuthURLOutputToResponse(output)
	lib.RespondJSON(w, http.StatusOK, resp)
}
