package http

import (
	"context"

	"Cryptoproject/internal/entities"
)

type CoinService interface {
	GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAggregateCoins(ctx context.Context, title []string, aggFunc string) ([]entities.Coin, error)
	StoreCoins(ctx context.Context, coins []entities.Coin) error
	GetCoinsList(ctx context.Context) ([]string, error)
}
