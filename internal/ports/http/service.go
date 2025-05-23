package http

import (
	"context"
	"net/http"

	"Cryptoproject/internal/entities"
)

type CoinService interface {
	GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAggregateCoins(ctx context.Context, title []string, aggFunc string) ([]entities.Coin, error)
	renderResponse(w http.ResponseWriter, status int, data interface{})
	renderError(w http.ResponseWriter, r *http.Request, err error)
}
