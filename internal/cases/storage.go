package cases

import (
	"context"

	"Cryptoproject/internal/entities"
)

type Storage interface {
	Store(ctx context.Context, coins []entities.Coin) error
	GetCoinsList(ctx context.Context) ([]string, error)
	GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAggregateCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
}
