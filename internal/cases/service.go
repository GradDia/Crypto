package cases

import (
	"context"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type Service interface {
	GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMaxRates(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMinRates(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvgRates(ctx context.Context, titles []string) ([]entities.Coin, error)
	ActualizeRates(ctx context.Context, opts ...Option) error
}
type ServiceImpl struct {
	storage        Storage
	cryptoProvider CryptoProvider
}

func NewService(storage Storage, cryptoProvider CryptoProvider) Service {
	return &ServiceImpl{
		storage:        storage,
		cryptoProvider: cryptoProvider,
	}
}

func (s *ServiceImpl) GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, entities.ErrInvalidParam
	}
	return s.storage.GetActualCoins(ctx, titles)
}

func (s *ServiceImpl) GetMaxRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, entities.ErrInvalidParam
	}
	return s.storage.GetAggregateCoins(ctx, titles)
}

func (s *ServiceImpl) GetMinRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, entities.ErrInvalidParam
	}
	return s.storage.GetAggregateCoins(ctx, titles)
}

func (s *ServiceImpl) GetAvgRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, entities.ErrInvalidParam
	}
	return s.storage.GetAggregateCoins(ctx, titles)
}

func (s *ServiceImpl) ActualizeRates(ctx context.Context, opts ...Option) error {
	config := &Options{}
	for _, opt := range opts {
		opt(config)
	}
	titles := config.Titles
	actualRates, err := s.cryptoProvider.GetActualRates(ctx, titles)
	if err != nil {
		return errors.Wrap(err, "failed to get actual rates from crypto provider")
	}
	for _, coin := range actualRates {
		if coin.CoinName == "" || coin.Price <= 0 {
			return errors.New("invalid coin data")
		}
		err := s.storage.Store(ctx, []entities.Coin{coin})
		if err != nil {
			return errors.Wrap(err, "failed to store update coin rate")
		}
	}
	return nil
}
