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
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "Cryptoproject/docs"
)

type Server struct {
	router      *chi.Mux
	httpServer  *http.Server
	coinService CoinService
	logger      *slog.Logger
}

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	if !w.wroteHeader {
		w.status = code
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
}

func NewServer(coinService CoinService, port string, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("component", "http-server"))

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешаем все origins для разработки
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кеширования preflight запросов
	}))

	s := &Server{
		router:      r,
		coinService: coinService,
		logger:      logger,
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
	}

	s.initRoutes()
	logger.Info("Server initialized", slog.String("port", port))
	return s
}

func (s *Server) initRoutes() {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriterWrapper{ResponseWriter: w}

			defer func() {
				s.logger.LogAttrs(
					r.Context(),
					slog.LevelInfo,
					"Request processed",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.status),
					slog.Duration("duration", time.Since(start)),
					slog.String("remote_addr", r.RemoteAddr),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	})

	s.router.Get("/api/v1/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger/doc.json"),
	))

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Post("/coins/actual", s.handleGetActualCoins)
		r.Post("/coins/aggregate/{aggFunc}", s.handleGetAggregateCoins)
	})

}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		slog.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
