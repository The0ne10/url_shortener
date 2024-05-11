package main

import (
	"fmt"
	"log/slog"
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
	log.Info("Logger initialized", slog.String("Server address", cfg.Address))
	log.Debug("Debug mode enabled")

	db, err := sqlite.NewConnect(cfg.Storage)
	if err != nil {
		fmt.Errorf("storage initialization failed: %s", err)
		return
	}

	id, err := db.SaveURL("https:google.com", "googlefuckoff")
	if err != nil {
		fmt.Errorf("failed to save URL: %s", err)
		return
	}
	fmt.Println(*id)

	//TODO: router (chi, chi-middleware, chi-render)

	//TODO: run server
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
