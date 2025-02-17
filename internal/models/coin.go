package models

import (
	"time"
)

// Представление сущности Coin для криптовалюты.

type Coin struct {
	ID            int       `json:"id"`
	CoinName      string    `json:"currency_name"`  // Название валюты
	Price         float64   `json:"price"`          // Текущая цена
	MinPrice      float64   `json:"min_price"`      // Минимальная цена за день
	MaxPrice      float64   `json:"max_price"`      // Максимальная цена за день
	ChangePercent float64   `json:"change_percent"` // Изменение в процентах за последний час
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

//NewCoin создание нового экземпляра Coin

func NewCoin(coinName string, price float64) *Coin {
	now := time.Now()
	return &Coin{
		CoinName:      coinName,
		Price:         price,
		MinPrice:      price, //Первое значение мин цены
		MaxPrice:      price, //Первое значение макс цены
		ChangePercent: 0,     //Первое изменение в процентах
		CreatedAt:     now,
		UpdatedAt:     now,
	}

}

// UpdatePrice обновление текущей, мин и макс цены.

func (c *Coin) UpdatePrice(newPrice float64) {
	c.Price = newPrice
	if newPrice < c.MinPrice || c.MinPrice == 0 {
		c.MinPrice = newPrice
	}
	if newPrice > c.MaxPrice {
		c.MaxPrice = newPrice
	}
	c.UpdatedAt = time.Now()
}

// CalculateChangePercent вычисляем изменение цены в процентах относительно предыдущей цены.

func (c *Coin) CalculateChangePercent(previousPrice float64) {
	if previousPrice == 0 {
		c.ChangePercent = 0
		return
	}
	c.ChangePercent = ((c.Price - previousPrice) / previousPrice) * 100

}
