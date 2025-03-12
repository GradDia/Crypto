package entities

import (
	"time"

	"github.com/pkg/errors"
)

type Coin struct {
	CoinName  string
	Price     float64
	CreatedAt time.Time
}

func NewCoin(coinName string, price float64) (*Coin, error) {
	if coinName == "" {
		return nil, errors.Wrap(ErrInvalidParam, "coin name not set")
	}
	if price <= 0 {
		return nil, errors.Wrap(ErrInvalidParam, "price must be greater then 0")
	}
	now := time.Now()
	return &Coin{
		CoinName:  coinName,
		Price:     price,
		CreatedAt: now,
	}, nil
}
