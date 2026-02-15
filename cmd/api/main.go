package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrielmatsan/checkin-gate/internal/config"
	eventshttp "github.com/gabrielmatsan/checkin-gate/internal/events/infra/http"
	identityhttp "github.com/gabrielmatsan/checkin-gate/internal/identity/infra/http"
	"github.com/gabrielmatsan/checkin-gate/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	_ "github.com/gabrielmatsan/checkin-gate/docs"
)

// @title           Checkin Gate API
// @version         1.0
// @description     API para sistema de check-in com autenticação OAuth Google

// @host      localhost:8080
// @BasePath  /

func main() {
	// Logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Redis
	redis, err := shared.NewRedis(cfg.RedisURL, logger)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer redis.Close()

	// Database
	db, err := shared.NewDatabase(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Health check
	if err := db.HealthCheck(); err != nil {
		logger.Fatal("database health check failed", zap.Error(err))
	}
	logger.Info("database connected")

	// Router
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// Routes
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Swagger JSON (necessário para o Scalar)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Scalar - documentação moderna da API
	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Checkin Gate API</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <script id="api-reference" data-url="/swagger/doc.json" data-configuration='{
        "theme": "purple",
        "layout": "modern",
        "darkMode": true,
        "hideDarkModeToggle": false,
        "searchHotKey": "k",
        "metaData": {
            "title": "Checkin Gate API",
            "description": "API para sistema de check-in com autenticação OAuth Google"
        },
        "hideDownloadButton": false
    }'></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`))
	})

	identityhttp.RegisterIdentityRoutes(router, db.DB, cfg)
	eventshttp.RegisterEventsRoutes(router, db.DB, cfg, logger)

	// Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	baseURL := fmt.Sprintf("http://localhost:%d", cfg.Port)

	fmt.Println("")
	fmt.Println("🚀 Checkin Gate API")
	fmt.Println("───────────────────────────────────────")
	fmt.Printf("🌐 Base:    %s\n", baseURL)
	fmt.Printf("❤️  Health:  %s/health\n", baseURL)
	fmt.Printf("📚 Docs:    %s/docs\n", baseURL)
	fmt.Printf("🔑 Auth:    %s/auth\n", baseURL)
	fmt.Println("───────────────────────────────────────")
	fmt.Println("")

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped")
}
