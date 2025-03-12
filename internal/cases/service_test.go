package cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/cases/mock_cases"
	"Cryptoproject/internal/entities"
)

func Test_GetLastRates_Success(t *testing.T) {
	t.Parallel()

	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для интерфейса Storage
	mockStorage := mock_cases.NewMockStorage(ctrl)

	// Определяем поведение мока
	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 28000},
		{CoinName: "ETH", Price: 1500},
	}
	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), titles).
		Return(expectedCoins, nil)

	// Создаем сервис с моком хранилища
	service := cases.NewService(mockStorage, nil)

	// Вызываем метод GetLastRates
	coins, err := service.GetLastRates(context.Background(), titles)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetLastRates_Error(t *testing.T) {
	t.Parallel()

	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для интерфейса Storage
	mockStorage := mock_cases.NewMockStorage(ctrl)

	// Определяем поведение мока
	titles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), titles).
		Return(nil, errors.New("storage error"))

	// Создаем сервис с моком хранилища
	service := cases.NewService(mockStorage, nil)

	// Вызываем метод GetLastRates
	coins, err := service.GetLastRates(context.Background(), titles)

	// Проверяем результат
	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.Contains(t, err.Error(), "storage error")
}
