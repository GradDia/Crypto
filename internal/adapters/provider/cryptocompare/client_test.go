package cryptocompare_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"Cryptoproject/internal/adapters/provider/cryptocompare"
	"Cryptoproject/internal/entities"
)

func Test_NewClient_Success(t *testing.T) {
	t.Parallel()

	apiKey := "test-api-key"
	client, err := cryptocompare.NewClient(apiKey)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.HttpClient)
}

func Test_NewClient_Error(t *testing.T) {
	t.Parallel()

	client, err := cryptocompare.NewClient("")
	require.Error(t, err)
	assert.True(t, errors.Is(err, entities.ErrInvalidParam))
	require.Nil(t, client)
	assert.Contains(t, err.Error(), "api key is required")
}

type mockTransport struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func Test_GetActualRates_Success(t *testing.T) {
	t.Parallel()

	client, err := cryptocompare.NewClient("test-api-key")
	require.NoError(t, err)

	// Setup mock transport
	client.HttpClient.Transport = &mockTransport{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "/data/pricemulti", req.URL.Path)
			require.Equal(t, "BTC,ETH", req.URL.Query().Get("fsyms"))
			require.Equal(t, "USD", req.URL.Query().Get("tsyms"))
			require.Equal(t, "Apikey test-api-key", req.Header.Get("Authorization"))

			response := map[string]map[string]float64{
				"BTC": {"USD": 50000.5},
				"ETH": {"USD": 3000.75},
			}

			w := httptest.NewRecorder()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return w.Result(), nil
		},
	}

	// Call method under test
	coins, err := client.GetActualRates(context.Background(), []string{"BTC", "ETH"})
	require.NoError(t, err)
	require.Len(t, coins, 2)

	assert.Equal(t, "BTC", coins[0].CoinName)
	assert.Equal(t, 50000.5, coins[0].Price)

	assert.Equal(t, "ETH", coins[1].CoinName)
	assert.Equal(t, 3000.75, coins[1].Price)
}

func Test_GetActualRates_ErrorCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		titles        []string
		mockResponse  interface{}
		mockStatus    int
		expectedError error
		errorContains string
	}{
		// ... другие тест-кейсы остаются без изменений ...
		{
			name:          "invalid json response",
			titles:        []string{"BTC"},
			mockResponse:  "invalid json",
			mockStatus:    http.StatusOK,
			expectedError: entities.ErrInternal,
			errorContains: "failed to decode response",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, err := cryptocompare.NewClient("test-api-key")
			require.NoError(t, err)

			client.HttpClient.Transport = &mockTransport{
				roundTrip: func(req *http.Request) (*http.Response, error) {
					w := httptest.NewRecorder()
					w.WriteHeader(tc.mockStatus)
					if tc.mockResponse != nil {
						switch v := tc.mockResponse.(type) {
						case string:
							w.Write([]byte(v))
						default:
							json.NewEncoder(w).Encode(v)
						}
					}
					return w.Result(), nil
				},
			}

			coins, err := client.GetActualRates(context.Background(), tc.titles)

			if tc.name == "invalid json response" {
				// Специальная обработка для случая с невалидным JSON
				require.Error(t, err)
				assert.Nil(t, coins)
				assert.Contains(t, err.Error(), tc.errorContains)

				// Проверяем, что это ошибка декодирования
				var jsonErr *json.SyntaxError
				if errors.As(err, &jsonErr) {
					// Это ожидаемая ошибка синтаксиса JSON
					return
				}

				// Дополнительная проверка для других типов ошибок декодирования
				var unmarshalErr *json.UnmarshalTypeError
				if errors.As(err, &unmarshalErr) {
					return
				}

				// Если ошибка не связана с декодированием JSON, проваливаем тест
				t.Errorf("expected JSON decoding error, got %T", err)
			} else {
				// Стандартная проверка для других случаев
				require.Error(t, err)
				assert.Nil(t, coins)
				assert.True(t, errors.Is(err, tc.expectedError),
					"expected error %v, got %v", tc.expectedError, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			}
		})
	}
}

func Test_WithPriceIn_Option(t *testing.T) {
	t.Parallel()

	client, err := cryptocompare.NewClient("test-api-key", cryptocompare.WithPriceIn("EUR"))
	require.NoError(t, err)

	client.HttpClient.Transport = &mockTransport{
		roundTrip: func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "EUR", req.URL.Query().Get("tsyms"))
			w := httptest.NewRecorder()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]map[string]float64{
				"BTC": {"EUR": 45000.0},
			})
			return w.Result(), nil
		},
	}

	coins, err := client.GetActualRates(context.Background(), []string{"BTC"})
	require.NoError(t, err)
	require.Len(t, coins, 1)
	assert.Equal(t, "BTC", coins[0].CoinName)
	assert.Equal(t, 45000.0, coins[0].Price)
}
