package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

const (
	AggFuncAVG = "AVG"
	AggFuncMAX = "MAX"
	AggFuncMin = "MIN"
)

func (s *Server) renderResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) renderError(w http.ResponseWriter, r *http.Request, err error) {

	// ErrorResponse represents error response structure
	// swagger:model ErrorResponse
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

	s.renderResponse(w, status, errorResponse{
		Error: err.Error(),
		Code:  status,
	})
}

// handleGetActualCoins godoc
// @Summary Get latest coin prices
// @Description Returns latest prices for requested coins
// @Tags coins
// @Accept json
// @Produce json
// @Param titles query string true "Comma-separated list of coin titles" Example("BTC,ETH")
// @Success 200 {array} entities.Coin
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coins/actual [post]
func (s *Server) handleGetActualCoins(w http.ResponseWriter, r *http.Request) {
	titlesParam := r.URL.Query().Get("titles")
	if titlesParam == "" {
		s.renderError(w, r, errors.Wrap(entities.ErrInvalidParam, "titles parameter is required"))
		return
	}

	titles := strings.Split(titlesParam, ",")
	if len(titles) == 0 {
		s.renderError(w, r, errors.Wrap(entities.ErrInvalidParam, "empty titles list"))
		return
	}

	coins, err := s.coinService.GetActualCoins(r.Context(), titles)
	if err != nil {
		s.renderError(w, r, errors.Wrap(err, "failed to get actual coins"))
		return
	}

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
// @Success 200 {array} entities.Coin
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coins/aggregate/{aggFunc} [post]
func (s *Server) handleGetAggregateCoins(w http.ResponseWriter, r *http.Request) {
	aggFunc := chi.URLParam(r, "aggFunc")

	validAggFuncs := map[string]bool{
		AggFuncAVG: true,
		AggFuncMAX: true,
		AggFuncMin: true,
	}

	if !validAggFuncs[aggFunc] {
		s.renderError(w, r, errors.Wrap(entities.ErrInvalidParam, "invali agg func, use: "+strings.Join([]string{AggFuncAVG, AggFuncMAX, AggFuncMin}, ",")))
		return
	}

	titlesParam := r.URL.Query().Get("titles")
	if titlesParam == "" {
		s.renderError(w, r, errors.Wrap(entities.ErrInvalidParam, "titles parameter is required"))
		return
	}

	titles := strings.Split(titlesParam, ",")
	if len(titles) == 0 {
		s.renderError(w, r, errors.Wrap(entities.ErrInvalidParam, "at last one title requred"))
	}

	coins, err := s.coinService.GetAggregateCoins(r.Context(), titles, titlesParam)
	if err != nil {
		s.renderError(w, r, errors.Wrap(err, "failed to get aggregate data"))
		return
	}

	s.renderResponse(w, http.StatusOK, coins)
}
