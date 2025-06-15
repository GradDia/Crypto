package cases

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type Service struct {
	storage        Storage
	cryptoProvider CryptoProvider
	logger         *slog.Logger
}

func NewService(storage Storage, cryptoProvider CryptoProvider, logger *slog.Logger) (*Service, error) {
	const op = "cases.NewService"
	if logger == nil {
		logger = slog.Default()
	}

	logger.Debug("Initializing new service",
		slog.Any("storage", storage != nil),
		slog.Any("cryptoProvider", cryptoProvider != nil))

	if storage == nil || storage == Storage(nil) {
		err := errors.Wrap(entities.ErrInvalidParam, "storage not set")
		logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}

	if cryptoProvider == nil || cryptoProvider == CryptoProvider(nil) {
		err := errors.Wrap(entities.ErrInvalidParam, "cryptoProvider not set")
		logger.Error(op, slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Service initialized Successfully")
	return &Service{
		storage:        storage,
		cryptoProvider: cryptoProvider,
		logger:         logger,
	}, nil
}

func (s *Service) GetLastRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	const op = "cases.GetLastRates"
	startTime := time.Now()
	logger := s.logger.With(slog.String("opertion", op))

	s.logger.Info("Processing request",
		slog.Any("titles_count", len(titles)),
		slog.Any("titles", titles))

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
		s.logger.Warn("Validation failed", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Debug("Getting fresh rates from provider")
	freshCoins, err := s.cryptoProvider.GetActualRates(ctx, titles)
	if err != nil {
		logger.Error("Failed to get fresh rates",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrap(err, "failed to get fresh rates")
	}

	logger.Debug("Shorting fresh rates",
		slog.Int("coins_count", len(freshCoins)))
	if err = s.storage.Store(ctx, freshCoins); err != nil {
		logger.Error("Failed to store fresh coins",
			slog.String("error", err.Error()),
			slog.Int("coins_count", len(freshCoins)))
		return nil, errors.Wrap(err, "failed to store fresh coins")
	}

	logger.Debug("Retrieving actual coins from storage")
	coins, err := s.storage.GetActualCoins(ctx, titles)
	if err != nil {
		logger.Error("Failed to get actual coins",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrap(err, "failed to get actual coins from storage")
	}

	logger.Info("Request processed successfully",
		slog.Int("coins_count", len(coins)),
		slog.Duration("duration", time.Since(startTime)))

	return coins, nil
}

func (s *Service) GetRatesWithAgg(ctx context.Context, titles []string, aggFuncTitle string) ([]entities.Coin, error) {
	const op = "cases.GetRatesWithAgg"
	startTime := time.Now()
	logger := s.logger.With(
		slog.String("opertion", op),
		slog.String("aggFuncTitle", aggFuncTitle))

	s.logger.Info("Processing aggregation request",
		slog.Any("titles_count", len(titles)))

	if len(titles) == 0 {
		err := errors.Wrap(entities.ErrInvalidParam, "titles list is empty")
		s.logger.Warn("Validation failed", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Debug("check existing titles")
	if err := s.checkExistingTitles(ctx, titles); err != nil {
		logger.Error("Failed to check existing titles",
			slog.String("error", err.Error()))
		return nil, errors.Wrap(err, "failed to check existing titles")
	}

	logger.Debug("Getting aggregated coins from storage")
	coins, err := s.storage.GetAggregateCoins(ctx, titles, aggFuncTitle)
	if err != nil {
		logger.Error("Failed to get aggregated coins",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrap(err, "failed to get aggregate coins from storage")
	}

	logger.Info("Aggregated complete",
		slog.Int("coins_count", len(coins)),
		slog.Duration("duration", time.Since(startTime)))

	return coins, nil
}

func (s *Service) ActualizeRates(ctx context.Context) error {
	const op = "cases.ActualizeRates"
	startTime := time.Now()
	logger := s.logger.With(slog.String("opertion", op))

	s.logger.Info("Starting rates actualization")

	logger.Debug("Getting coins list from storage")
	allTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		s.logger.Error("failed to get coins list",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return errors.Wrap(err, "actualizeRates get coins list")
	}

	s.logger.Debug("Retrieved coins list",
		slog.Int("count", len(allTitles)))

	logger.Debug("Getting actual rates from provider")
	updatedCoins, err := s.cryptoProvider.GetActualRates(ctx, allTitles)
	if err != nil {
		s.logger.Error("failed to get actual rates",
			slog.String("error", err.Error()))
		return errors.Wrap(err, "actualizeRates get actual rates")
	}

	s.logger.Debug("Shorting updated rates",
		slog.Int("coins_count", len(updatedCoins)))
	if err = s.storage.Store(ctx, updatedCoins); err != nil {
		s.logger.Error("failed to store actual rates",
			slog.String("error", err.Error()),
			slog.Int("coins_count", len(updatedCoins)))
		return errors.Wrap(err, "actualizeRates store")
	}

	s.logger.Info("Rates actualization completed successfully",
		slog.Int("coins_updated", len(updatedCoins)),
		slog.Duration("duration", time.Since(startTime)))
	return nil
}

func (s *Service) checkExistingTitles(ctx context.Context, requestTitles []string) error {
	const op = "cases.checkExistingTitles"
	logger := s.logger.With(slog.String("op", op))

	logger.Debug("Checking existing titles",
		slog.Int("requst_titles_count", len(requestTitles)))

	existingTitles, err := s.storage.GetCoinsList(ctx)
	if err != nil {
		logger.Error("Failed to get coins list",
			slog.String("error", err.Error()))
		return errors.Wrap(err, "failed to get coins list from storage")
	}

	missingTitles := findMissingTitles(requestTitles, existingTitles)
	if len(missingTitles) == 0 {
		logger.Debug("All titles exist in storage")
		return nil
	}

	logger.Info("Found missing titles",
		slog.Int("missing_count", len(missingTitles)),
		slog.Any("missing_titles", missingTitles))

	logger.Debug("Getting rates for missing titles")
	newCoins, err := s.cryptoProvider.GetActualRates(ctx, missingTitles)
	if err != nil {
		logger.Error("Failed to get rates for missing titles",
			slog.String("error", err.Error()),
			slog.Any("missing_titles", missingTitles))
		return errors.Wrap(err, "failed to get actual rates for missing titles")
	}

	logger.Debug("Storing missing titles",
		slog.Int("new_coins_count", len(newCoins)))
	if err = s.storage.Store(ctx, newCoins); err != nil {
		logger.Error("Failed to store new coins",
			slog.String("error", err.Error()),
			slog.Int("coins_count", len(newCoins)))
		return errors.Wrap(err, "failed to store new coins in storage")
	}

	logger.Info("Missing titles processed successfully",
		slog.Int("added_count", len(newCoins)))

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
