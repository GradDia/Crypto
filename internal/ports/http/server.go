package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router      *chi.Mux
	httpServer  *http.Server
	coinService CoinService
}

func NewServer(coinService CoinService, port string) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

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
		r.Get("/coins/actual", s.handleGetActualCoins)
		r.Get("/coins/aggregate")
		r.Get("/coins")
		r.Get("/coins/list")
	})
	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
