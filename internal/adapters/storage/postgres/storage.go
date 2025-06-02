package postgres

import (
	"Cryptoproject/internal/cases"
	"Cryptoproject/internal/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"strings"
)

var (
	_ cases.Storage = (*Storage)(nil)
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(connectionString string) (*Storage, error) {
	if connectionString == "" {
		return nil, errors.Wrap(entities.ErrInvalidParam, "missing connection string")
	}

	con, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "failed to create connection pool, err: %s", err.Error())
	}
	return &Storage{db: con}, nil
}

func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	names := make([]string, 0, len(coins))
	prices := make([]float64, 0, len(coins))

	for _, coin := range coins {
		names = append(names, coin.CoinName)
		prices = append(prices, coin.Price)
	}

	_, err := s.db.Exec(ctx, `
		INSERT INTO coins (coin_name, price)
		SELECT unnest($1::text[]), unnest($2::decimal[])
		`, names, prices)

	if err != nil {
		return errors.Wrapf(entities.ErrInternal, "failed to store coins: %s", err.Error())
	}
	return nil
}

func (s *Storage) GetCoinsList(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT DISTINCT coin_name FROM coins`)
	if err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
		}
		coins = append(coins, name)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}
	return coins, nil
}

func (s *Storage) GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
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

	rows, err := s.db.Query(ctx, query, params...)

	if err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}
	defer rows.Close()

	coins := make([]entities.Coin, 0, len(titles))
	for rows.Next() {
		var coin entities.Coin
		if err = rows.Scan(&coin.CoinName, &coin.Price, &coin.CreatedAt); err != nil {
			return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
		}
		coins = append(coins, coin)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "failed to query coins list")
	}
	return coins, nil
}

func (s *Storage) GetAggregateCoins(ctx context.Context, titles []string, aggFunc string) ([]entities.Coin, error) {

	if len(titles) == 0 {
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
		return nil, errors.Wrap(entities.ErrInternal, "unsupported aggregate function: %s (allowed: AVG, MAX, MIN)")
	}

	rows, err := s.db.Query(ctx, aggQuery, titles)

	if err != nil {
		return nil, errors.Wrapf(entities.ErrInternal, "aggregate query failed, err: %v", err)
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err = rows.Scan(&coin.CoinName, &coin.Price); err != nil {
			return nil, errors.Wrap(entities.ErrInternal, "failed to scan coin data")
		}
		coins = append(coins, coin)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "error during rows iteration")
	}
	return coins, nil
}
