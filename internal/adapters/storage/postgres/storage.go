package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"math"
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
	const maxRetries = 5
	const retryDelay = 2 * time.Second

	if logger == nil {
		logger = slog.Default()
	}
	logger = logger.With(
		slog.String("component", "postgres-storage"),
		slog.String("op", op),
	)

	logger.Debug("Initializing PostgreSQL storage",
		slog.String("connection_string", maskPassword(connectionString))) // Маскируем пароль в логах

	if connectionString == "" {
		err := errors.Wrap(entities.ErrInvalidParam, "missing connection string")
		logger.Error("Invalid parameter")
		return nil, err
	}

	var pool *pgxpool.Pool
	var err error
	startTime := time.Now()

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Info("Connecting to database",
			slog.Int("attempt", attempt),
			slog.String("host", extractHost(connectionString)))

		pool, err = pgxpool.New(context.Background(), connectionString)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err = pool.Ping(ctx); err == nil {
				logger.Info("Database connection established",
					slog.Duration("duration", time.Since(startTime)))
				return &Storage{db: pool, logger: logger}, nil
			}
			pool.Close()
		}

		if attempt < maxRetries {
			delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			logger.Warn("Connection failed, retrying",
				slog.String("error", err.Error()),
				slog.Duration("delay", delay))

			time.Sleep(delay)
			continue
		}
	}

	logger.Error("Failed to connect after all retries",
		slog.String("error", err.Error()),
		slog.Duration("total_duration", time.Since(startTime)))
	return nil, errors.Wrap(err, "failed to connect after retries")
}

func maskPassword(connStr string) string {
	if strings.Contains(connStr, "@") {
		parts := strings.Split(connStr, "@")
		if len(parts) > 0 {
			return parts[0] + "@*****"
		}
	}
	return connStr
}

func extractHost(connStr string) string {
	if strings.Contains(connStr, "@") {
		parts := strings.Split(connStr, "@")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return connStr
}

func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	const op = "postgres.Store"
	logger := s.logger.With(
		slog.String("op", op),
		slog.Int("coins_count", len(coins)),
	)

	// Простая вставка без проверки уникальности
	_, err := s.db.Exec(ctx, `
        INSERT INTO coins (coin_name, price, created_at)
        SELECT unnest($1::text[]), unnest($2::decimal[]), now()
    `,
		getCoinNames(coins),  // []string названий монет
		getCoinPrices(coins), // []float64 цен
	)

	if err != nil {
		logger.Error("Insert failed", slog.String("error", err.Error()))
		return errors.Wrap(err, "failed to insert coins")
	}

	logger.Info("Coins stored successfully")
	return nil
}

// Вспомогательные функции
func getCoinNames(coins []entities.Coin) []string {
	names := make([]string, len(coins))
	for i, c := range coins {
		names[i] = c.CoinName
	}
	return names
}

func getCoinPrices(coins []entities.Coin) []float64 {
	prices := make([]float64, len(coins))
	for i, c := range coins {
		prices[i] = c.Price
	}
	return prices
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
