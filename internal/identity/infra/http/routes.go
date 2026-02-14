package http

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	authenticatewithgoogle "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/authenticate_with_google"
	getgoogleauthurl "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/get_google_auth_url"
	refreshtoken "github.com/gabrielmatsan/checkin-gate/internal/identity/application/usecase/refresh_token"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/http/handler"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/persistence"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra/service"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterIdentityRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	googleProvider := lib.NewGoogleOAuthProvider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	jwtService := service.NewJWTService(cfg.JWTSecret)
	userRepo := persistence.NewPostgresUserRepository(db)
	sessionRepo := persistence.NewPostgresSessionRepository(db)

	getGoogleAuthURL := getgoogleauthurl.NewUseCase(googleProvider)
	authenticateWithGoogle := authenticatewithgoogle.NewUseCase(googleProvider, jwtService, userRepo, sessionRepo)
	refreshToken := refreshtoken.NewUseCase(jwtService, userRepo, sessionRepo)

	// Create individual handlers
	googleCallbackHandler := handler.NewGoogleCallbackHandler(authenticateWithGoogle)
	refreshTokenHandler := handler.NewRefreshTokenHandler(refreshToken)
	getGoogleAuthURLHandler := handler.NewGetGoogleAuthURLHandler(getGoogleAuthURL)

	r.Route("/auth", func(r chi.Router) {
		r.Get("/google/url", getGoogleAuthURLHandler.Handle)
		r.Get("/google/callback", googleCallbackHandler.Handle)

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))
			r.Post("/refresh", refreshTokenHandler.Handle)
		})
	})
}

// GetJWTService returns a new JWT service for use by other modules
func GetJWTService(cfg *config.Config) *service.JWTService {
	return service.NewJWTService(cfg.JWTSecret)
}
