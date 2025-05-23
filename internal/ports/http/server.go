// @title Cryptocurrency API
// @version 1.0
// @description API for cryptocurrency data management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@cryptoproject.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
package http

import (
	"context"
	"net/http"

	_ "Cryptoproject/docs"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
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
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(
			"doc.json"),
		))
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
