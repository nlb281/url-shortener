package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
)

const (
	envLocal = "local"
	envDev = "development"
)

func main() {
	config := config.MustLoad()

	fmt.Println(config)

	log := initLogger(config.Env)

	log.Info("starting url-shortener", slog.String("env", config.Env))
	
	// TODO: init storage

	// TODO: init router

	// TODO: run server
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}