package cases

import (
	"context"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type Service struct {
	storage        Storage
	cryptoProvider CryptoProvider
}

func NewService(storage Storage, cryptoProvider CryptoProvider) (*Service, error) {
	if storage == nil || storage == Storage(nil) {
		return nil, errors.Wrap(entities.ErrInvalidParam, "storage not set")
	}

	if cryptoProvider == nil || cryptoProvider == CryptoProvider(nil) {
		return nil, errors.Wrap(entities.ErrInvalidParam, "cryptoProvider not set")
	}

	return &Service{
		storage:        storage,
		cryptoProvider: cryptoProvider,
	}, nil
}

func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	//	err := s.checkExistingTitles(ctx, titles)
	//	if err != nil {
	//		return nil, errors.Wrap(entities.ErrInvalidParam, "failed to check existing titles")
	//	}

	coins, err := s.storage.GetActualCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "failed to get actual coins from storage")
	}
	return coins, nil
}

func (s *Service) GetMaxRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	//err := s.checkExistingTitles(ctx, titles)
	//if err != nil {
	//	return nil, errors.Wrap(entities.ErrInvalidParam, "failed to check existing titles")
	//}

	coins, err := s.storage.GetAggregateCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "failed to get aggregate coins from storage")
	}

	var maxRates []entities.Coin
	for _, coin := range coins {
		if coin.Price > 0 {
			maxRates = append(maxRates, entities.Coin{
				CoinName: coin.CoinName,
				Price:    coin.Price,
			})
		}
	}
	return maxRates, nil
}

func (s *Service) GetMinRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	//err := s.checkExistingTitles(ctx, titles)
	//if err != nil {
	//	return nil, errors.Wrap(entities.ErrInvalidParam, "failed to check existing titles")
	//}

	coins, err := s.storage.GetAggregateCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "failed to get aggregate coins from storage")
	}

	var minRates []entities.Coin
	for _, coin := range coins {
		if coin.Price > 0 {
			minRates = append(minRates, entities.Coin{
				CoinName: coin.CoinName,
				Price:    coin.Price,
			})
		}
	}
	return minRates, nil
}

func (s *Service) GetAvgRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	//err := s.checkExistingTitles(ctx, titles)
	//if err != nil {
	//	return nil, errors.Wrap(entities.ErrInvalidParam, "failed to check existing titles")
	//}

	coins, err := s.storage.GetAggregateCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInvalidParam, "failed to get aggregate coins from storage")
	}

	var avgRates []entities.Coin
	for _, coin := range coins {
		if coin.Price > 0 {
			avgRates = append(avgRates, entities.Coin{
				CoinName: coin.CoinName,
				Price:    coin.Price,
			})
		}
	}
	return avgRates, nil
}

func (s *Service) ActualizeRates(ctx context.Context, opts ...Option) error {
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

func (s *Service) CheckExistingTitles(ctx context.Context, requestTitles []string) error {
	existingTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get coins list from storage")
	}

	missingTitles := findMissingTitles(requestTitles, existingTitles)
	if len(missingTitles) == 0 {
		return nil
	}

	newCoins, err := s.cryptoProvider.GetActualRates(ctx, missingTitles)
	if err != nil {
		return errors.Wrap(err, "failed to get actual rates for missing titles")
	}

	for _, coin := range newCoins {
		err := s.storage.Store(ctx, []entities.Coin{coin})
		if err != nil {
			return errors.Wrap(err, "failed to store new coins in storage")
		}
	}
	return nil
}

func findMissingTitles(requestTitles, existingTitles []string) []string {
	existingSet := make(map[string]struct{})
	for _, title := range existingTitles {
		existingSet[title] = struct{}{}
	}

	var missingTitles []string
	for _, title := range requestTitles {
		if _, exists := existingSet[title]; !exists {
			missingTitles = append(missingTitles, title)
		}
	}
	return missingTitles
}
