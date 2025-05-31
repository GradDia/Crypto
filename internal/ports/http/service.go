package http

import (
	"Cryptoproject/internal/entities"
	"context"
)

type CoinService interface {
	GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetRatesWithAgg(ctx context.Context, titles []string, aggFunc string) ([]entities.Coin, error)
}
