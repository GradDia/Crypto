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

	if err := s.checkExistingTitles(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "failed to check existing titles")
	}

	coins, err := s.storage.GetActualCoins(ctx, titles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get actual coins from storage")
	}
	return coins, nil
}

func (s *Service) GetRatesWithAgg(ctx context.Context, titles []string, aggFuncTitle string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
	}

	if err := s.checkExistingTitles(ctx, titles); err != nil {
		return nil, errors.Wrap(err, "failed to check existing titles")
	}

	coins, err := s.storage.GetAggregateCoins(ctx, titles, aggFuncTitle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get aggregate coins from storage")
	}
	return coins, nil
}

func (s *Service) ActualizeRates(ctx context.Context) error {
	allTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		return errors.Wrap(err, "actualizeRates get coins list")
	}

	updatedCoins, err := s.cryptoProvider.GetActualRates(ctx, allTitles)
	if err != nil {
		return errors.Wrap(err, "actualizeRates get actual rates")
	}

	if err = s.storage.Store(ctx, updatedCoins); err != nil {
		return errors.Wrap(err, "actualizeRates store")
	}
	return nil
}

func (s *Service) checkExistingTitles(ctx context.Context, requestTitles []string) error {
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

	if err = s.storage.Store(ctx, newCoins); err != nil {
		return errors.Wrap(err, "failed to store new coins in storage")
	}
	return nil
}

func findMissingTitles(requestTitles, existingTitles []string) []string {
	existingSet := make(map[string]struct{}, len(existingTitles))
	for _, title := range existingTitles {
		existingSet[title] = struct{}{}
	}

	missingTitles := make([]string, 0, len(requestTitles))
	for _, title := range requestTitles {
		if _, exists := existingSet[title]; !exists {
			missingTitles = append(missingTitles, title)
		}
	}
	return missingTitles
}
