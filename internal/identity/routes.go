package identity

import (
	"github.com/gabrielmatsan/checkin-gate/internal/config"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/handler"
	"github.com/gabrielmatsan/checkin-gate/internal/identity/infra"
	usecases "github.com/gabrielmatsan/checkin-gate/internal/identity/usecases"
	"github.com/gabrielmatsan/checkin-gate/internal/shared/lib"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	googleProvider := lib.NewGoogleOAuthProvider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	jwtService := infra.NewJWTService(cfg.JWTSecret)
	userRepo := infra.NewPostgresUserRepository(db)
	sessionRepo := infra.NewPostgresSessionRepository(db)

	authenticateWithGoogle := usecases.NewAuthenticateWithGoogleUseCase(googleProvider, jwtService, userRepo, sessionRepo)
	refreshToken := usecases.NewRefreshTokenUseCase(jwtService, userRepo, sessionRepo)

	authHandler := handler.NewAuthHandler(authenticateWithGoogle, refreshToken)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/google/callback", authHandler.GoogleCallback)
		r.Post("/refresh", authHandler.Refresh)
	})
}
