package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url_shortener/internal/config"
	"url_shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "develop"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()

	//TODO: logger
	log := setupLogger(cfg.Env)

	log = log.With("Env", cfg.Env)
	log.Info("Logger initialized")
	log.Debug("Debug mode enabled")

	_, err := sqlite.NewConnect(cfg.Storage)
	if err != nil {
		log.Error("storage initialization failed: %s", err)
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.URLFormat)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("Server initialized", slog.String("Server address", srv.Addr))

	err = srv.ListenAndServe()
	if err != nil {
		return
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
