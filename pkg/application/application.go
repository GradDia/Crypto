package application

import (
	"context"
	"log"

	"Cryptoproject/internal/adapters/provider/cryptocompare"
	"Cryptoproject/internal/adapters/storage/postgres"
	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/ports/http"
)

type App struct {
	httpServer *http.Server
}

func NewApp() (*App, error) {
	storage, err := postgres.NewStorage("postgres://user:password@localhost:5232/coins?sslmode=disable")
	if err != nil {
		return nil, err
	}

	cryptoProvider, err := cryptocompare.NewClient("your_cryptocompare_api_key")
	if err != nil {
		return nil, err
	}

	service, err := cases.NewService(storage, cryptoProvider)
	if err != nil {
		return nil, err
	}

	httpServer := http.NewServer(service, "8080")

	return &App{
		httpServer: httpServer,
	}, nil
}

func (a *App) Run() error {
	log.Println("Starting server on :8080")
	return a.httpServer.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.httpServer != nil {
		return a.httpServer.Stop(ctx)
	}
	return nil
}
