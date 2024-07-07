package main

import (
	"golang-websocket-chat/internal/config"
	httpserver "golang-websocket-chat/internal/httpServer"
	"golang-websocket-chat/internal/storage/postgres"
	sl "golang-websocket-chat/lib/logger/slog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := sl.New(cfg.Env)

	log.Info("staring chat backend", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	rs := httpserver.New(storage.Rooms(), storage.Users(), storage.Messages(), log)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      rs,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done

	log.Info("stopping server")

	log.Info("server stopped")
}
