package dto

import "time"

// CoinResponse DTO для ответа API (актуальные данные)
// swagger:model CoinResponse
type CoinResponse struct {
	CoinName  string    `json:"coin_name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"` // omit epty
}

// AggregateCoinResponse DTO для агрегированных данных
// swagger:model AggregateCoinResponse
type AggregateCoinResponse struct {
	CoinName string  `json:"coin_name"`
	Price    float64 `json:"price"` // AVG, MAX или MIN
}
