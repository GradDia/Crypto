// Package http implements HTTP server for cryptocurrency API
//
// @title Cryptocurrency API
// @version 1.0
// @description API for cryptocurrency data management
// @termsOfService http://swagger.io/terms/
//
// @contact.name @GradDia
// @contact.email not support
// @license.name for free
// @license.url for free
//
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package http

import (
	_ "Cryptoproject/docs"

	"context"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router      *chi.Mux
	httpServer  *http.Server
	coinService CoinService
}

func NewServer(coinService CoinService, port string) *Server {
	r := chi.NewRouter()

	s := &Server{
		router:      r,
		coinService: coinService,
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
	}

	s.initRoutes()
	return s
}

func (s *Server) initRoutes() {
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Указываем явный путь
	))

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Post("/coins/actual", s.handleGetActualCoins)
		r.Post("/coins/aggregate", s.handleGetAggregateCoins)
	})

}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
