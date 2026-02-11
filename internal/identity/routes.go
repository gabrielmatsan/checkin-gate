package identity

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/handler"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra"
	usecases "github.com/gabrielmatsan/checkin-gate/internal/identity/usecases"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	googleProvider := lib.NewGoogleOAuthProvider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	jwtService := infra.NewJWTService(cfg.JWTSecret)
	userRepo := infra.NewPostgresUserRepository(db)
	sessionRepo := infra.NewPostgresSessionRepository(db)

	getGoogleAuthURL := usecases.NewGetGoogleAuthURLUseCase(googleProvider)

	authenticateWithGoogle := usecases.NewAuthenticateWithGoogleUseCase(googleProvider, jwtService, userRepo, sessionRepo)

	refreshToken := usecases.NewRefreshTokenUseCase(jwtService, userRepo, sessionRepo)

	authHandler := handler.NewAuthHandler(getGoogleAuthURL, authenticateWithGoogle, refreshToken)

	r.Route("/auth", func(r chi.Router) {

		r.Get("/google/url", authHandler.GetGoogleAuthURL)
		r.Get("/google/callback", authHandler.GoogleCallback)

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(middleware.NewValidateTokenFunc(jwtService.ExtractClaims)))
			r.Post("/refresh", authHandler.Refresh)
			//r.Delete("/logout", )
		})
	})
}
