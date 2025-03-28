package cases_test

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/cases/testdata"
	"Cryptoproject/internal/entities"
)

func Test_GetLastRates_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 28000},
		{CoinName: "ETH", Price: 1500},
	}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return([]string{"BTC", "ETH"}, nil)

	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), titles).
		Return(expectedCoins, nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetLastRates(context.Background(), titles)
	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetLastRates_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	titles := []string{"BTC", "ETH"}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return([]string{"BTC", "ETH"}, nil)

	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), titles).
		Return(nil, entities.ErrInvalidParam)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetLastRates(context.Background(), titles)
	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.ErrorIs(t, err, entities.ErrInvalidParam)
}

func Test_GetRatesWithAgg_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 30000},
		{CoinName: "ETH", Price: 2000},
	}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return([]string{"BTC", "ETH"}, nil)

	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles, "max").
		Return(expectedCoins, nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetRatesWithAgg(context.Background(), titles, "max")
	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetRatesWithAgg_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	titles := []string{"BTC", "ETH"}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return([]string{"BTC", "ETH"}, nil)

	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles, "max").
		Return(nil, entities.ErrInvalidParam)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetRatesWithAgg(context.Background(), titles, "max")
	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.ErrorIs(t, err, entities.ErrInvalidParam)
}

func Test_ActualizeRates_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	// Определяем поведение мока Storage
	allTitles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(allTitles, nil)

	// Определяем поведение мока CryptoProvider
	actualRates := []entities.Coin{
		{CoinName: "BTC", Price: 29000},
		{CoinName: "ETH", Price: 1600},
	}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), allTitles).
		Return(actualRates, nil)

	// Определяем поведение мока Storage для сохранения новых монет
	mockStorage.EXPECT().
		Store(gomock.Any(), actualRates).
		Return(nil)

	// Создаем сервис с моками
	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	// Вызываем метод ActualizeRates
	err = service.ActualizeRates(context.Background())
	assert.NoError(t, err)
}

func Test_ActualizeRates_CryptoProviderError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	allTitles := []string{"BTC", "ETH"}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(allTitles, nil)

	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), allTitles).
		Return(nil, errors.New("crypto provider error"))

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	err = service.ActualizeRates(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "crypto provider error")
}

func Test_GetLastRates_AllCoinsExist(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	// Определяем поведение моков
	requestTitles := []string{"BTC", "ETH"}
	existingTitles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 28000},
		{CoinName: "ETH", Price: 1500},
	}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(existingTitles, nil)

	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), requestTitles).
		Return(expectedCoins, nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetLastRates(context.Background(), requestTitles)
	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetLastRates_AddMissingCoins(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	// Определяем поведение моков
	requestTitles := []string{"BTC", "ETH", "LTC"}
	existingTitles := []string{"BTC", "ETH"}
	missingTitles := []string{"LTC"}

	newCoins := []entities.Coin{
		{CoinName: "LTC", Price: 50},
	}

	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(existingTitles, nil)

	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), missingTitles).
		Return(newCoins, nil)

	for _, coin := range newCoins {
		mockStorage.EXPECT().
			Store(gomock.Any(), []entities.Coin{coin}).
			Return(nil)
	}

	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), requestTitles).
		Return(append([]entities.Coin{
			{CoinName: "BTC", Price: 28000},
			{CoinName: "ETH", Price: 1500},
		}, newCoins...), nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetLastRates(context.Background(), requestTitles)
	assert.NoError(t, err)
	assert.Equal(t, append([]entities.Coin{
		{CoinName: "BTC", Price: 28000},
		{CoinName: "ETH", Price: 1500},
	}, newCoins...), coins)
}
