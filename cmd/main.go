package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/exp/slog"

	"github.com/Nicktsim/kodetest/config"
	"github.com/Nicktsim/kodetest/handlers/create"
	"github.com/Nicktsim/kodetest/handlers/get"
	"github.com/Nicktsim/kodetest/handlers/users/login"
	"github.com/Nicktsim/kodetest/handlers/users/register"
	"github.com/Nicktsim/kodetest/logger/sl"
	mwLogger "github.com/Nicktsim/kodetest/middleware"
	"github.com/Nicktsim/kodetest/storage/psql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("starting kode-test-project",
			slog.String("address", cfg.Address))
	log.Debug("debug messages are enabled")
	storageParams := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	storage, err := psql.NewStorage(storageParams)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/signup", register.SignUp(log, storage))

    router.Post("/signin", login.SignIn(log, storage))

    router.Post("/create", create.NewNote(log, storage))

    router.Get("/notes", get.GetUserNotes(log, storage))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}
	defer storage.Close()

	log.Info("server stopped")
}
