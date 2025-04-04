package cryptocompare_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Cryptoproject/internal/adapters/provider/cryptocompare"
	"Cryptoproject/internal/entities"
)

func Test_GetActualRates_Success(t *testing.T) {
	t.Parallel()

	// Создаем фиктивный HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовки и параметры запроса
		assert.Equal(t, "Apikey test-api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "/data/pricemulti?fsyms=BTC,ETH&tsyms=USD", r.URL.Path+"?"+r.URL.RawQuery)

		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"BTC":{"USD":28000},"ETH":{"USD":1500}}`))
	}))
	defer server.Close()

	// Преобразуем server.URL (string) в *url.URL
	proxyURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	// Создаем клиент с фиктивным URL
	client, err := cryptocompare.NewClient("test-api-key")
	require.NoError(t, err)
	client.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	// Вызываем метод GetActualRates
	titles := []string{"BTC", "ETH"}
	expectedCoins := []entities.Coin{
		{CoinName: "BTC", Price: 28000},
		{CoinName: "ETH", Price: 1500},
	}

	coins, err := client.GetActualRates(context.Background(), titles)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, expectedCoins, coins)
}

func Test_GetActualRates_EmptyTitles(t *testing.T) {
	t.Parallel()

	// Создаем клиент
	client, err := cryptocompare.NewClient("test-api-key")
	require.NoError(t, err)

	// Вызываем метод с пустым списком тайтлов
	coins, err := client.GetActualRates(context.Background(), []string{})

	// Проверяем ошибку
	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.Contains(t, err.Error(), "titles list is empty")
}

func Test_GetActualRates_APIError(t *testing.T) {
	t.Parallel()

	// Создаем фиктивный HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
	}))
	defer server.Close()

	// Преобразуем server.URL (string) в *url.URL
	proxyURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	// Создаем клиент с фиктивным URL
	client, err := cryptocompare.NewClient("test-api-key")
	require.NoError(t, err)
	client.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	// Вызываем метод GetActualRates
	titles := []string{"BTC", "ETH"}
	coins, err := client.GetActualRates(context.Background(), titles)

	// Проверяем ошибку
	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

func Test_GetActualRates_DecodeError(t *testing.T) {
	t.Parallel()

	// Создаем фиктивный HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid-json-response`))
	}))
	defer server.Close()

	// Преобразуем server.URL (string) в *url.URL
	proxyURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	// Создаем клиент с фиктивным URL
	client, err := cryptocompare.NewClient("test-api-key")
	require.NoError(t, err)
	client.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	// Вызываем метод GetActualRates
	titles := []string{"BTC", "ETH"}
	coins, err := client.GetActualRates(context.Background(), titles)

	// Проверяем ошибку
	assert.Error(t, err)
	assert.Nil(t, coins)
	assert.Contains(t, err.Error(), "failed to decode response")
}
