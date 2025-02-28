package cases_test

import (
	"context"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"Cryptoproject/internal/cases"
)

func Test_CreateCoin_Succes(t *testing.T) {
	t.Parallel()

	service := cases.NewService()

	coinName := "BTC"
	price := 2700.0
	coin, err := service.CreateCoin(context.Background(), coinName, price)

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	if coin == nil {
		t.Fatal("Expected coin to be created, but got nil")
	}
	if coin.CoinName != coinName {
		t.Errorf("Expected coin name to be '%s' but got: %s", coinName, coin.CoinName)
	}
	if coin.Price != price {
		t.Errorf("Expected price to be %f, but got: %f", price, coin.Price)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func Test_CreateCoin_ValidateError(t *testing.T) {
	t.Parallel()

	// Создаем сервис
	service := cases.NewService()

	// Попытка создать монету с пустым именем
	coinName := ""
	price := 1.0
	coin, err := service.CreateCoin(context.Background(), coinName, price)

	// Проверяем результат
	if !errors.Is(err, cases.ErrInvalidData) {
		t.Errorf("Expected error '%v', but got: %v", cases.ErrInvalidData, err)
	}
	if coin != nil {
		t.Fatalf("Expected coin to be nil, but got: %+v", coin)
	}
	if !contains(err.Error(), "The name of the coin cannot be empty") {
		t.Errorf("Expected error message to contain 'The name of the coin cannot be empty', but got: %v", err)
	}

	// Попытка создать монету с нулевой ценой
	coinName = "BTC"
	price = 0.0
	coin, err = service.CreateCoin(context.Background(), coinName, price)

	// Проверяем результат
	if !errors.Is(err, cases.ErrInvalidData) {
		t.Errorf("Expected error '%v', but got: %v", cases.ErrInvalidData, err)
	}
	if coin != nil {
		t.Fatalf("Expected coin to be nil, but got: %+v", coin)
	}
	if !contains(err.Error(), "The price must be more than 0") {
		t.Errorf("Expected error message to contain 'The price must be more than 0', but got: %v", err)
	}
}

func Test_GetCoinByName_NotFound(t *testing.T) {
	t.Parallel()

	service := cases.NewService()

	coinName := "BTC"
	coin, err := service.GetCoinByName(context.Background(), coinName)

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
	if coin != nil {
		t.Fatalf("Expected coin to be nil, but got: %+v", coin)
	}
	if err.Error() != "Coin not found" {
		t.Errorf("Expected error 'Coin not found', But got: %v", err)
	}
}

func Test_UpdateCoinPrice_NotFound(t *testing.T) {
	t.Parallel()

	// Создаем сервис
	service := cases.NewService()

	// Вызываем метод UpdateCoinPrice
	coinName := "BTC"
	newPrice := 28000.0
	coin, err := service.UpdateCoinPrice(context.Background(), coinName, newPrice)

	// Проверяем результат
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
	if coin != nil {
		t.Fatalf("Expected coin to be nil, but got: %+v", coin)
	}
	if err.Error() != "Coin not found" {
		t.Errorf("Expected error 'Coin not found', but got: %v", err)
	}
}
