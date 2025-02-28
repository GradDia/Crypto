package cases

import (
	"context"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type Service interface {
	CreateCoin(ctx context.Context, coinName string, price float64) (*entities.Coin, error)
	GetCoinByName(ctx context.Context, coinName string) (*entities.Coin, error)
	UpdateCoinPrice(ctx context.Context, coinName string, newPrice float64) (*entities.Coin, error)
}

func NewService() Service {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (s *serviceImpl) CreateCoin(ctx context.Context, coinName string, price float64) (*entities.Coin, error) {
	if coinName == "" {
		return nil, errors.Wrap(ErrInvalidData, "The name of the coin cannot be empty")
	}
	if price <= 0 {
		return nil, errors.Wrap(ErrInvalidData, "The price must be more than 0")
	}
	return entities.NewCoin(coinName, price)
}

func (s *serviceImpl) GetCoinByName(ctx context.Context, coinName string) (*entities.Coin, error) {
	return nil, errors.New("Coin not found")
}

func (s *serviceImpl) UpdateCoinPrice(ctx context.Context, coinName string, newPrice float64) (*entities.Coin, error) {
	return nil, errors.New("Coin not found")
}
