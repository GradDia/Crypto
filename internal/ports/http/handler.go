package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
	"Cryptoproject/pkg/dto"
)

const (
	AggFuncAVG = "AVG"
	AggFuncMAX = "MAX"
	AggFuncMin = "MIN"
)

func (s *Server) renderResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Добавьте эту строку
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("Failed to render response",
			slog.Int("status", status),
			slog.String("error", err.Error()))
	}
}

func (s *Server) renderError(w http.ResponseWriter, r *http.Request, err error) {
	const op = "http.renderError"
	logger := s.logger.With(
		slog.String("op", op),
		slog.String("path", r.URL.Path),
	)

	type errorResponse struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
	}

	var status int
	switch {
	case errors.Is(err, entities.ErrInvalidParam):
		status = http.StatusBadRequest
	case errors.Is(err, entities.ErrNotFound):
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	logger.Warn("Rendering error response",
		slog.Int("status", status),
		slog.String("error", err.Error()))

	s.renderResponse(w, status, errorResponse{
		Error: err.Error(),
		Code:  status,
	})
}

// handleGetActualCoins godoc
// @Summary Get latest coin prices
// @Description Returns latest prices for requested coins. Tickers must be comma-separated without spaces (e.g. "BTC,ETH")
// @Tags coins
// @Accept json
// @Produce json
// @Param titles query string true "Comma-separated list of coin titles without spaces" Example("BTC,ETH")
// @Success 200 {array} dto.CoinResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 500 {object} dto.ErrorResponseDto
// @Router /api/v1/coins/actual [post]
func (s *Server) handleGetActualCoins(w http.ResponseWriter, r *http.Request) {
	const op = "http.handleGetActualCoins"
	startTime := time.Now()
	logger := s.logger.With(
		slog.String("op", op),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	titlesParam := r.URL.Query().Get("titles")
	if titlesParam == "" {
		err := errors.Wrap(entities.ErrInvalidParam, "titles parameter is required")
		logger.Warn("Validation failed", slog.String("error", err.Error()))
		s.renderError(w, r, err)
		return
	}

	titlesParam = strings.ReplaceAll(titlesParam, " ", "")
	titles := strings.Split(titlesParam, ",")
	logger = logger.With(slog.Int("titles_count", len(titles)))

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "empty titles list")
		logger.Warn("Validation failed", slog.String("error", err.Error()))
		s.renderError(w, r, err)
		return
	}

	logger.Debug("Processing request", slog.Any("titles", titles))
	coins, err := s.coinService.GetLastRates(r.Context(), titles)
	if err != nil {
		logger.Error("Service call failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		s.renderError(w, r, errors.Wrap(err, "failed to get actual coins"))
		return
	}
	//convert coins to response dto
	response := make([]dto.CoinResponse, 0, len(coins))
	for _, coin := range coins {
		response = append(response, dto.CoinResponse{
			CoinName:  coin.CoinName,
			Price:     coin.Price,
			CreatedAt: coin.CreatedAt,
		})
	}

	logger.Info("Request processed successfully",
		slog.Int("coins_count", len(response)),
		slog.Duration("duration", time.Since(startTime)))

	s.renderResponse(w, http.StatusOK, coins)
}

// handleGetAggregateCoins godoc
// @Summary Get aggregated coin data
// @Description Returns aggregated data (AVG/MAX/MIN) for requested coins
// @Tags coins
// @Accept json
// @Produce json
// @Param aggFunc path string true "Aggregation function (AVG, MAX, MIN)" Enums(AVG, MAX, MIN)
// @Param titles query string true "Comma-separated list of coin titles" Example("BTC,ETH")
// @Success 200 {array} dto.AggregateCoinResponse
// @Failure 400 {object} dto.ErrorResponseDto
// @Failure 500 {object} dto.ErrorResponseDto
// @Router /api/v1/coins/aggregate/{aggFunc} [post]
func (s *Server) handleGetAggregateCoins(w http.ResponseWriter, r *http.Request) {
	const op = "http.handleGetAggregateCoins"
	startTime := time.Now()
	logger := s.logger.With(
		slog.String("op", op),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	aggFunc := chi.URLParam(r, "aggFunc")
	logger = logger.With(slog.String("agg_func", aggFunc))

	validAggFuncs := map[string]struct{}{
		AggFuncAVG: {},
		AggFuncMAX: {},
		AggFuncMin: {},
	}

	if _, ok := validAggFuncs[aggFunc]; !ok {
		err := errors.Wrapf(entities.ErrInvalidParam, "invalid agg func: %s", aggFunc)
		logger.Warn("Validation failed", slog.String("error", err.Error()))
		s.renderError(w, r, err)
		return
	}

	titlesParam := r.URL.Query().Get("titles")
	if titlesParam == "" {
		err := errors.Wrap(entities.ErrInvalidParam, "titles parameter is required")
		logger.Warn("Validation failed", slog.String("error", err.Error()))
		s.renderError(w, r, err)
		return
	}

	titlesParam = strings.ReplaceAll(titlesParam, " ", "")

	titles := strings.Split(titlesParam, ",")
	logger = logger.With(slog.Int("titles_count", len(titles)))

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "at least one title required")
		logger.Warn("Validation failed", slog.String("error", err.Error()))
		s.renderError(w, r, err)
	}

	logger.Debug("Processing aggregation request",
		slog.Any("titles", titles),
		slog.String("agg_func", aggFunc))

	aggregateData, err := s.coinService.GetRatesWithAgg(r.Context(), titles, aggFunc)
	if err != nil {
		logger.Error("Service call failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		s.renderError(w, r, errors.Wrap(err, "failed to get aggregate data"))
		return
	}

	response := make([]dto.AggregateCoinResponse, 0, len(aggregateData))
	for _, data := range aggregateData {
		response = append(response, dto.AggregateCoinResponse{
			CoinName: data.CoinName,
			Price:    data.Price,
		})
	}

	logger.Info("Aggregation request processed",
		slog.Int("coins_count", len(response)),
		slog.Duration("duration", time.Since(startTime)))

	s.renderResponse(w, http.StatusOK, response)
}
