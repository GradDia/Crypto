package postgres

import (
	"Cryptoproject/internal/cases"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
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
	names := make([]string, len(coins))
	prices := make([]float64, len(coins))

	for i, coin := range coins {
		names[i] = coin.CoinName
		prices[i] = coin.Price
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

	params := make([]interface{}, 0, len(titles))

	for _, title := range titles {
		params = append(params, title)
	}

	rows, err := s.db.Query(ctx, `
        SELECT DISTINCT ON (coin_name) coin_name, price, created_at
        FROM coins
        WHERE coin_name = ANY($1)
        ORDER BY coin_name, created_at DESC
		`, params...)

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

func (s *Storage) GetAggregateCoins(ctx context.Context, titles []string, aggFuncTitle string) ([]entities.Coin, error) {

	if len(titles) == 0 {
		return []entities.Coin{}, nil
	}

	var query string
	switch aggFuncTitle {
	case "AVG":
		query = "AVG(price)"
	case "MAX":
		query = "MAX(price)"
	case "MIN":
		query = "MIN(price)"
	default:
		return nil, errors.Wrap(entities.ErrInternal, "unsupported aggregate function: %s (allowed: AVG, MAX, MIN)")
	}

	rows, err := s.db.Query(ctx, `
        SELECT 
            coin_name, 
            `+query+` as price,
            MAX(created_at) as created_at
        FROM coins
        WHERE coin_name = ANY($1)
        GROUP BY coin_name
        ORDER BY coin_name
    `, titles)

	if err != nil {
		return nil, errors.Wrap(entities.ErrInternal, "aggregate query failed")
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
