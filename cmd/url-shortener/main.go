package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/storage/db"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	
	storage, err := db.New(config.StoragePath)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}


	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Start server", slog.String("Address", config.Address))

	srv := &http.Server{
		Addr: 				config.Address,
		Handler: 			router,
		ReadTimeout: 	config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout: 	config.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
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