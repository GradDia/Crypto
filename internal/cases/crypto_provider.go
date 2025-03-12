package cases

import (
	"context"

	"Cryptoproject/internal/entities"
)

type CryptoProvider interface {
	GetActualRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}
