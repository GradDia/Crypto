package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/entities"
)

var (
	_ cases.Storage = (*Storage)(nil)
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewStorage(connectionString string, logger *slog.Logger) (*Storage, error) {
	const op = "postgres.NewStorage"
	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(slog.String("component", "postgres-storage"))

	logger.Debug("Initializing PostgreSQL storage")

	if connectionString == "" {
		err := errors.Wrap(entities.ErrInvalidParam, "missing connection string")
		logger.Error("Invalid parameter", slog.String("error", err.Error()))
		return nil, err
	}

	startTime := time.Now()
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.Error("Connection pool creation failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrap(err, "failed to create connection pool")
	}

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		logger.Error("Database ping failed",
			slog.String("error", err.Error()))
		return nil, errors.Wrap(err, "database ping failed")
	}

	logger.Info("Storage initialized successfully",
		slog.Duration("duration", time.Since(startTime)))

	return &Storage{db: pool, logger: logger}, nil
}

func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	const op = "postgres.Store"
	logger := s.logger.With(
		slog.String("op", op),
		slog.Int("coins_count", len(coins)),
	)
	startTime := time.Now()

	if len(coins) == 0 {
		logger.Debug("Empty coins list provided")
		return nil
	}

	names := make([]string, 0, len(coins))
	prices := make([]float64, 0, len(coins))
	for _, coin := range coins {
		names = append(names, coin.CoinName)
		prices = append(prices, coin.Price)
	}

	logger.Debug("Executing batch insert")
	_, err := s.db.Exec(ctx, `
		INSERT INTO coins (coin_name, price)
		SELECT unnest($1::text[]), unnest($2::decimal[])
	`, names, prices)

	if err != nil {
		logger.Error("Insert operation failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return errors.Wrap(err, "failed to store coins")
	}

	logger.Info("Coins stored successfully",
		slog.Duration("duration", time.Since(startTime)))
	return nil
}

func (s *Storage) GetCoinsList(ctx context.Context) ([]string, error) {
	const op = "postgres.GetCoinsList"
	logger := s.logger.With(
		slog.String("op", op),
	)
	startTime := time.Now()

	logger.Debug("Querying distinct coins")
	rows, err := s.db.Query(ctx, `SELECT DISTINCT coin_name FROM coins`)
	if err != nil {
		logger.Error("Query failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			logger.Error("Row scan failed",
				slog.String("error", err.Error()))
			return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
		}
		coins = append(coins, name)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Rows iteration failed",
			slog.String("error", err.Error()))
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}

	logger.Info("Coins list retrieved",
		slog.Int("count", len(coins)),
		slog.Duration("duration", time.Since(startTime)))
	return coins, nil
}

func (s *Storage) GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	const op = "postgres.GetActualCoins"
	logger := s.logger.With(
		slog.String("op", op),
		slog.Int("titles_count", len(titles)),
	)
	startTime := time.Now()

	if len(titles) == 0 {
		logger.Debug("Empty titles list provided")
		return []entities.Coin{}, nil
	}

	placeholders := make([]string, len(titles))
	params := make([]interface{}, len(titles))
	for i, title := range titles {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		params[i] = title
	}

	query := fmt.Sprintf(`
    SELECT DISTINCT ON (coin_name) coin_name, price, created_at
    FROM coins
    WHERE coin_name IN(%s)
    ORDER BY coin_name, created_at DESC
`, strings.Join(placeholders, ","))

	logger.Debug("Executing query",
		slog.String("query", query))
	rows, err := s.db.Query(ctx, query, params...)

	if err != nil {
		logger.Error("Query failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}
	defer rows.Close()

	coins := make([]entities.Coin, 0, len(titles))
	for rows.Next() {
		var coin entities.Coin
		if err = rows.Scan(&coin.CoinName, &coin.Price, &coin.CreatedAt); err != nil {
			logger.Error("Row scan failed",
				slog.String("error", err.Error()))
			return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
		}
		coins = append(coins, coin)
	}
	if err = rows.Err(); err != nil {
		logger.Error("Rows iteration failed",
			slog.String("error", err.Error()))
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}

	logger.Info("Actual coins retrieved",
		slog.Int("count", len(coins)),
		slog.Duration("duration", time.Since(startTime)))
	return coins, nil
}

func (s *Storage) GetAggregateCoins(ctx context.Context, titles []string, aggFunc string) ([]entities.Coin, error) {
	const op = "postgres.GetAggregateCoins"
	logger := s.logger.With(
		slog.String("op", op),
		slog.String("agg_func", aggFunc),
		slog.Int("titles_count", len(titles)),
	)
	startTime := time.Now()

	if len(titles) == 0 {
		logger.Debug("Empty titles list provided")
		return []entities.Coin{}, nil
	}

	query := `
        SELECT 
            coin_name, 
            %s(price) as price
        FROM coins
        WHERE coin_name = ANY($1)
        GROUP BY coin_name
    `

	var aggQuery string
	switch strings.ToUpper(aggFunc) {
	case "AVG":
		aggQuery = fmt.Sprintf(query, "AVG")
	case "MAX":
		aggQuery = fmt.Sprintf(query, "MAX")
	case "MIN":
		aggQuery = fmt.Sprintf(query, "MIN")
	default:
		err := errors.Wrap(entities.ErrInternal, "unsupported aggregate function: %s (allowed: AVG, MAX, MIN)")
		logger.Error("Invalid aggregation function",
			slog.String("error", err.Error()))
		return nil, err
	}

	logger.Debug("Executing aggregate query",
		slog.String("query", aggQuery))
	rows, err := s.db.Query(ctx, aggQuery, titles)

	if err != nil {
		logger.Error("Aggregate query failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(startTime)))
		return nil, errors.Wrapf(entities.ErrInternal, "aggregate query failed, err: %v", err)
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err = rows.Scan(&coin.CoinName, &coin.Price); err != nil {
			logger.Error("Row scan failed",
				slog.String("error", err.Error()))
			return nil, errors.Wrap(entities.ErrInternal, "failed to scan coin data")
		}
		coins = append(coins, coin)
	}
	if err = rows.Err(); err != nil {
		logger.Error("Rows iteration failed",
			slog.String("error", err.Error()))
		return nil, errors.Wrap(entities.ErrInternal, "error during rows iteration")
	}

	logger.Info("Aggregate data retrieved",
		slog.Int("count", len(coins)),
		slog.Duration("duration", time.Since(startTime)))
	return coins, nil
}
