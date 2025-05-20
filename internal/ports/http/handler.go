package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

func renderResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	type errorResponse struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
	}

	var status int
	switch {
	case errors.Is(err, entities.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, entities.ErrNotFound):
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	renderResponse(w, r, status, errorResponse{
		Error: err.Error(),
		Code:  status,
	})
}

func (s *Server) handleGetActualCoins(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Titles []string `json:"titles"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "invalid request body"))
		return
	}

	if len(req.Titles) == 0 {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "empty titles list"))
		return
	}

	coins, err := s.coinService.GetActualCoins(r.Context(), req.Titles)
	if err != nil {
		renderError(w, r, errors.Wrap(err, "failed to get actual coins"))
		return
	}

	renderResponse(w, r, http.StatusOK, coins)
}

func (s *Server) handleGetAggregateCoins(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Titles  []string `json:"titles"`
		AggFunc string   `json:"agg_func"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "invalid request body"))
		return
	}

	if len(req.Titles) == 0 {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "empty titles list"))
		return
	}

	switch req.AggFunc {
	case "AVG", "MAX", "MIN":
	default:
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "invalid agg func"))
		return
	}

	coins, err := s.coinService.GetAggregateCoins(r.Context(), req.Titles, req.AggFunc)
	if err != nil {
		renderError(w, r, errors.Wrap(err, "failed to get aggregate data"))
		return
	}

	renderResponse(w, r, http.StatusOK, coins)
}

func (s *Server) handleStoreCoins(w http.ResponseWriter, r *http.Request) {
	var coins []entities.Coin

	if err := json.NewDecoder(r.Body).Decode(&coins); err != nil {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "invalid request body"))
		return
	}

	if len(coins) == 0 {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "empty coins list"))
		return
	}

	if err := s.coinService.StoreCoins(r.Context(), coins); err != nil {
		renderError(w, r, errors.Wrap(entities.ErrInvalidInput, "failed to store coins"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleGetCoinsList(w http.ResponseWriter, r *http.Request) {
	list, err := s.coinService.GetCoinsList(r.Context())
	if err != nil {
		renderError(w, r, errors.Wrap(err, "failet to get coins list"))
		return
	}

	renderResponse(w, r, http.StatusOK, list)
}
