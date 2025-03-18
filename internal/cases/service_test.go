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
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 28000},
		{CoinName: "ETH", Price: 1500},
	}
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
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetActualCoins(gomock.Any(), titles).
		Return(nil, entities.ErrInvalidParam) // Возвращаем базовую ошибку

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err) // Проверяем, что сервис создан без ошибок

	coins, err := service.GetLastRates(context.Background(), titles)

	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.ErrorIs(t, err, entities.ErrInvalidParam) // Проверяем тип ошибки
}

func Test_GetMaxRates_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 30000},
		{CoinName: "ETH", Price: 2000},
	}
	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles).
		Return(expectedCoins, nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetMaxRates(context.Background(), titles)

	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}
func Test_GetMaxRates_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles).
		Return(nil, entities.ErrInvalidParam)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetMaxRates(context.Background(), titles)

	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.ErrorIs(t, err, entities.ErrInvalidParam)
}

func Test_GetMinRates_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 27000},
		{CoinName: "ETH", Price: 1400},
	}
	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles).
		Return(expectedCoins, nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetMinRates(context.Background(), titles)

	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetMinRates_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles).
		Return(nil, entities.ErrInvalidParam)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetMinRates(context.Background(), titles)

	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.ErrorIs(t, err, entities.ErrInvalidParam)
}

func Test_GetAvgRates_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 28500},
		{CoinName: "ETH", Price: 1700},
	}
	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles).
		Return(expectedCoins, nil)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetAvgRates(context.Background(), titles)

	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetAvgRates_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	titles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetAggregateCoins(gomock.Any(), titles).
		Return(nil, entities.ErrInvalidParam)

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	coins, err := service.GetAvgRates(context.Background(), titles)

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

	titles := []string{"BTC", "ETH"}
	actualRates := []entities.Coin{
		{CoinName: "BTC", Price: 29000},
		{CoinName: "ETH", Price: 1600},
	}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), titles).
		Return(actualRates, nil)

	for _, coin := range actualRates {
		mockStorage.EXPECT().
			Store(gomock.Any(), []entities.Coin{coin}).
			Return(nil)
	}

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	err = service.ActualizeRates(context.Background(), cases.WithTitles(titles))

	assert.NoError(t, err)
}

func Test_ActualizeRates_CryptoProviderError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	titles := []string{"BTC", "ETH"}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), titles).
		Return(nil, errors.New("crypto provider error"))

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	err = service.ActualizeRates(context.Background(), cases.WithTitles(titles))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "crypto provider error")
}

func Test_ActualizeRates_StorageError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)

	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	titles := []string{"BTC", "ETH"}
	actualRates := []entities.Coin{
		{CoinName: "BTC", Price: 29000},
		{CoinName: "ETH", Price: 1600},
	}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), titles).
		Return(actualRates, nil)

	for i, coin := range actualRates {
		if i == 0 { // Первый вызов успешен
			mockStorage.EXPECT().
				Store(gomock.Any(), []entities.Coin{coin}).
				Return(nil)
		} else { // Второй вызов возвращает ошибку
			mockStorage.EXPECT().
				Store(gomock.Any(), []entities.Coin{coin}).
				Return(errors.New("storage error"))
		}
	}

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	err = service.ActualizeRates(context.Background(), cases.WithTitles(titles))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage error")
}

func Test_CheckExistingTitles_Success(t *testing.T) {
	t.Parallel()

	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для интерфейса Storage
	mockStorage := testdata.NewMockStorage(ctrl)

	// Определяем поведение мока Storage
	requestTitles := []string{"BTC", "ETH"}
	existingTitles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(existingTitles, nil)

	// Создаем мок для интерфейса CryptoProvider
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), gomock.Any()).
		Return([]entities.Coin{}, nil).AnyTimes()

	// Создаем сервис с моками
	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	// Вызываем метод checkExistingTitles
	err = service.CheckExistingTitles(context.Background(), requestTitles)

	// Проверяем результат
	assert.NoError(t, err)
}

func Test_CheckExistingTitles_AddMissingCoins(t *testing.T) {
	t.Parallel()

	// Создаем контроллер для моков
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для интерфейса Storage
	mockStorage := testdata.NewMockStorage(ctrl)

	// Создаем мок для интерфейса CryptoProvider
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	// Определяем поведение мока Storage
	requestTitles := []string{"BTC", "ETH", "LTC"}
	existingTitles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(existingTitles, nil)

	// Определяем поведение мока CryptoProvider
	missingTitles := []string{"LTC"}
	newCoins := []entities.Coin{
		{CoinName: "LTC", Price: 50},
	}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), missingTitles).
		Return(newCoins, nil)

	// Определяем поведение мока Storage для сохранения новых монет
	for _, coin := range newCoins {
		mockStorage.EXPECT().
			Store(gomock.Any(), []entities.Coin{coin}).
			Return(nil)
	}

	// Создаем сервис с моками
	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	// Вызываем метод checkExistingTitles
	err = service.CheckExistingTitles(context.Background(), requestTitles)

	// Проверяем результат
	assert.NoError(t, err)
}

func Test_CheckExistingTitles_CryptoProviderError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	requestTitles := []string{"BTC", "ETH", "LTC"}
	existingTitles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(existingTitles, nil)

	missingTitles := []string{"LTC"}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), missingTitles).
		Return(nil, errors.New("crypto provider error"))

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	err = service.CheckExistingTitles(context.Background(), requestTitles)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "crypto provider error")
}

func Test_CheckExistingTitles_StorageError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := testdata.NewMockStorage(ctrl)
	mockCryptoProvider := testdata.NewMockCryptoProvider(ctrl)

	requestTitles := []string{"BTC", "ETH", "LTC"}
	existingTitles := []string{"BTC", "ETH"}
	mockStorage.EXPECT().
		GetCoinsList(gomock.Any()).
		Return(existingTitles, nil)

	missingTitles := []string{"LTC"}
	newCoins := []entities.Coin{
		{CoinName: "LTC", Price: 50},
	}
	mockCryptoProvider.EXPECT().
		GetActualRates(gomock.Any(), missingTitles).
		Return(newCoins, nil)

	for _, coin := range newCoins {
		mockStorage.EXPECT().
			Store(gomock.Any(), []entities.Coin{coin}).
			Return(errors.New("storage error"))
	}

	service, err := cases.NewService(mockStorage, mockCryptoProvider)
	require.NoError(t, err)

	err = service.CheckExistingTitles(context.Background(), requestTitles)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage error")
}
