package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "Cryptoproject/docs"

	"Cryptoproject/pkg/application"
)

func main() {
	app := application.NewApp()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	log.Println("Server started on :8080")

	<-done
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

	log.Println("Server Stopped")
}
