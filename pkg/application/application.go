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

func NewApp() *App {
	storage, err := postgres.NewStorage("postgres://user:password@localhost:5432/coins?sslmode=disable")
	if err != nil {
		panic(err)
	}

	cryptoProvider, err := cryptocompare.NewClient("d994a06fb570c9237da009fcb028d0094662333e9f6d7e231707198442174ac5")
	if err != nil {
		panic(err)
	}

	service, err := cases.NewService(storage, cryptoProvider)
	if err != nil {
		panic(err)
	}

	httpServer := http.NewServer(service, "8080")

	return &App{
		httpServer: httpServer,
	}
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
