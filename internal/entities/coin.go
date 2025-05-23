package entities

import (
	"time"

	"github.com/pkg/errors"
)

// Coin represents cryptocurrency data
// swagger:model Coin
type Coin struct {
	CoinName  string    `json:"coin_name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
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
