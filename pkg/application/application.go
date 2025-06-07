package application

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"Cryptoproject/internal/adapters/provider/cryptocompare"
	"Cryptoproject/internal/adapters/storage/postgres"
	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/ports/http"
)

type App struct {
	httpServer *http.Server
	cron       *cron.Cron
	service    *cases.Service
}

func NewApp() *App {
	storage, err := postgres.NewStorage("postgres://user:password@localhost:5432/coins?sslmode=disable")
	if err != nil {
		panic(err)
	}

	cryptoProvider, err := cryptocompare.NewClient("your_api_key")
	if err != nil {
		panic(err)
	}

	service, err := cases.NewService(storage, cryptoProvider)
	if err != nil {
		panic(err)
	}

	httpServer := http.NewServer(service, "8080")

	app := &App{
		httpServer: httpServer,
		service:    service,
		cron:       cron.New(),
	}

	app.setupCron()
	return app
}

func (a *App) setupCron() {
	_, err := a.cron.AddFunc("*/10 * * * *", a.updateCoinData)
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	go func() {
		log.Println("Starting cron scheduler...")
		a.cron.Start()
	}()
}

func (a *App) updateCoinData() {
	startTime := time.Now()
	log.Println("[Cron] Starting coins data update...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.service.ActualizeRates(ctx); err != nil {
		log.Printf("[Cron] Update failed: %v", err)
		return
	}

	log.Printf("[Cron] Update completed in %v", time.Since(startTime))
}

func (a *App) Run() error {
	log.Println("Starting server on :8080")
	return a.httpServer.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down application...")

	if a.cron != nil {
		log.Println("Stopping cron scheduler...")
		cronCtx := a.cron.Stop()
		<-cronCtx.Done()
	}

	if a.httpServer != nil {
		log.Println("Stopping HTTP server...")
		if err := a.httpServer.Stop(ctx); err != nil {
			return err
		}
	}

	log.Println("Application shutdown completed")
	return nil
}
