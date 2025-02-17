package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestNewCoin проверка создания нового Coin через конструктор

func TestNewCoin(t *testing.T) {
	coin := NewCoin("BTC", 27000)

	// Проверка, что все поля инициализированы
	assert.Equal(t, "BTC", coin.CoinName, "Название монеты будет 'BTC'")
	assert.Equal(t, 27000.0, coin.Price, "Начальная цена должна быть 100")
	assert.Equal(t, 27000.0, coin.MinPrice, "Минимальная цена должна быть равна начальной цене")
	assert.Equal(t, 27000.0, coin.MaxPrice, "Максимальная цена должны быть равно начальной цене")
	assert.Equal(t, 0.0, coin.ChangePercent, "Изначальный процент изменения будет 0")
	assert.True(t, coin.CreatedAt.Before(time.Now()), "Создание должно быть в прошлом")
	assert.True(t, coin.UpdatedAt.Before(time.Now()), "Обновление должно быть в прошлом")
}

// TestUpdatePrice проверяем обновление цены

func TestUpdatePrice(t *testing.T) {
	coin := NewCoin("BTC", 27000)

	// Обновляем цену на более высокую
	coin.UpdatePrice(28000)
	assert.Equal(t, 28000.0, coin.Price, "Цена обновилась до 28000")
	assert.Equal(t, 27000.0, coin.MinPrice, "Минимальная цена должна быть 27000")
	assert.Equal(t, 28000.0, coin.MaxPrice, "Максимальная цена должна быть 28000")

	//Обновление цены на меньшую
	coin.UpdatePrice(26000)
	assert.Equal(t, 26000.0, coin.Price, "Цена обновилась до 26000")
	assert.Equal(t, 26000.0, coin.MinPrice, "Минимальная цена измениться до 26000")
	assert.Equal(t, 28000.0, coin.MaxPrice, "Максимальная цена останеться на 28000")
}

// TestCalculateChangePercent проверка вычисления изменений цены в процентах.
func TestCalculateChangePercent(t *testing.T) {
	coin := NewCoin("BTC", 27000)

	// Вычисляем изменение относительно предыдущей цены
	coin.CalculateChangePercent(25000)
	assert.Equal(t, 8.0, coin.ChangePercent, "Процент изменения цены составил 8.0%")

	// Проверяем случай, когда предыдущая цена равна текущей
	coin.CalculateChangePercent(27000)
	assert.Equal(t, 0.0, coin.ChangePercent, "Процент изменения цены составил 0%, цена не изменилась")

	// Проверяем, когда предыдущая цена равна 0
	coin.CalculateChangePercent(0)
	assert.Equal(t, 0.0, coin.ChangePercent, "Процент изменения будет 0.0%, если предыдущая цена равна 0")
}
