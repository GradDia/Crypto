package http

import (
	"context"
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
