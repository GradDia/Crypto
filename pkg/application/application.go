package application

import (
	"context"
	"log/slog"
	"os"
	"time"

	cronJob "github.com/robfig/cron/v3"

	"Cryptoproject/internal/adapters/provider/cryptocompare"
	"Cryptoproject/internal/adapters/storage/postgres"
	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/ports/http"
)

type App struct {
	httpServer *http.Server
	cron       *cronJob.Cron
	service    *cases.Service
	logger     *slog.Logger
}

func NewApp() *App {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	logger.Info("Initializing application")

	// Передаем логгер в NewStorage
	storage, err := postgres.NewStorage("postgres://user:password@localhost:5432/coins?sslmode=disable", logger)
	if err != nil {
		logger.Error("Failed to initialize storage", slog.String("error", err.Error()))
		panic(err)
	}

	// Передаем логгер в NewClient
	cryptoProvider, err := cryptocompare.NewClient("your_api_key", logger)
	if err != nil {
		logger.Error("Failed to initialize crypto provider", slog.String("error", err.Error()))
		panic(err)
	}

	// Передаем логгер в NewService
	service, err := cases.NewService(storage, cryptoProvider, logger)
	if err != nil {
		logger.Error("Failed to initialize service", slog.String("error", err.Error()))
		panic(err)
	}

	// Передаем логгер в NewServer
	httpServer := http.NewServer(service, "8080", logger)

	app := &App{
		httpServer: httpServer,
		service:    service,
		cron: cronJob.New(cronJob.WithLogger(
			cronJob.VerbosePrintfLogger(slog.NewLogLogger(logger.Handler(), slog.LevelDebug)),
		)),
		logger: logger,
	}

	app.setupCron()
	logger.Info("Application initialized successfully")
	return app
}

func (a *App) setupCron() {
	const jobName = "coin_data_update"

	_, err := a.cron.AddFunc("*/10 * * * *", func() {
		ctx := context.Background()
		startTime := time.Now()
		logger := a.logger.With(
			slog.String("job", jobName),
			slog.Time("start_time", startTime),
		)

		logger.Info("Starting cron job execution")

		if err := a.service.ActualizeRates(ctx); err != nil {
			logger.Error("Cron job failed",
				slog.String("error", err.Error()),
				slog.Duration("duration", time.Since(startTime)))
			return
		}

		logger.Info("Cron job completed successfully",
			slog.Duration("duration", time.Since(startTime)))
	})

	if err != nil {
		a.logger.Error("Failed to schedule cron job",
			slog.String("job", jobName),
			slog.String("error", err.Error()))
		panic(err)
	}

	go func() {
		a.logger.Info("Starting cron scheduler")
		a.cron.Start()
	}()
}

func (a *App) updateCoinData() {
	ctx := context.Background()
	startTime := time.Now()
	logger := a.logger.With(
		slog.String("op", "updateCoinData"),
		slog.Time("start_time", startTime),
	)

	logger.Info("Starting coins data update")

	if err := a.service.ActualizeRates(ctx); err != nil {
		logger.Error("Update failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return
	}

	logger.Info("Update completed",
		slog.Duration("duration", time.Since(startTime)))
}

func (a *App) Run() error {
	a.logger.Info("Starting HTTP server", slog.String("port", "8080"))
	return a.httpServer.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	// Остановка cron
	if a.cron != nil {
		a.logger.Info("Stopping cron scheduler")
		cronCtx := a.cron.Stop()
		select {
		case <-cronCtx.Done():
			a.logger.Info("Cron scheduler stopped gracefully")
		case <-ctx.Done():
			a.logger.Warn("Forced cron scheduler shutdown")
		}
	}

	// Остановка HTTP сервера
	if a.httpServer != nil {
		a.logger.Info("Stopping HTTP server")
		if err := a.httpServer.Stop(ctx); err != nil {
			a.logger.Error("Failed to stop HTTP server",
				slog.String("error", err.Error()))
			return err
		}
		a.logger.Info("HTTP server stopped")
	}

	a.logger.Info("Application shutdown completed")
	return nil
}
