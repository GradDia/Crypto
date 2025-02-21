package entities_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"Cryptoproject/internal/entities"
)

func Test_NewCoin_Success(t *testing.T) {
	t.Parallel()

	coinName := "test name"
	price := 0.1
	coin, err := entities.NewCoin(coinName, price)
	require.NoError(t, err)
	require.Equal(t, coinName, coin.CoinName)
	require.Equal(t, price, coin.Price)
}

func Test_NewCoin_ValidateError(t *testing.T) {
	t.Parallel()

	coinName := ""
	price := 1.0
	coin, err := entities.NewCoin(coinName, price)
	require.ErrorIs(t, err, entities.ErrInvalidParam)
	require.Nil(t, coin)
	require.Contains(t, err.Error(), "coin name not set")

	coinName = "name"
	price = 0.0
	coin, err = entities.NewCoin(coinName, price)
	require.ErrorIs(t, err, entities.ErrInvalidParam)
	require.Nil(t, coin)
	require.Contains(t, err.Error(), "price must be greater then 0")
}
