package postgres

import (
	"Cryptoproject/internal/entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Storage struct {
	db *pgxpool.Pool
}

func newStorage(connectionString string) (*Storage, error) {
	con, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection pool")
	}
	if err := con.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return &Storage{db: con}, nil
}

func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {
	batch := pgx.Batch{}

	for _, coin := range coins {
		batch.Queue(`
			INSERT INTO coins (coin_name, price, created_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (coin_name) DO UPDATE
			SET price = EXCLUDED.price, created_at = NOW()
			`, coin.CoinName, coin.Price)
	}
	results := s.db.SendBatch(ctx, &batch)
	defer results.Close()

	for range coins {
		_, err := results.Exec()
		if err != nil {
			return errors.Wrap(err, "failed to batch query")
		}
	}
	return nil
}

func (s *Storage) GetCoinsList(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, "SELECT coin_name FROM coins")
	if err != nil {
		return nil, errors.Wrap(err, "failed to query coins list")
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, errors.Wrap(err, "failed to query coins list")
		}
		coins = append(coins, name)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to query coins list")
	}
	return coins, nil
}

func (s *Storage) GetActualCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	rows, err := s.db.Query(ctx, `
		SELECT coin_name, price, created_at
		FROM coins
		WHERE coin_name = ANY($1)
		ORDER BY created_at DESC
		`, titles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query coins list")
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err := rows.Scan(&coin.CoinName, &coin.Price, &coin.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to query coins list")
		}
		coins = append(coins, coin)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to query coins list")
	}
	return coins, nil
}

func (s *Storage) AggregateCoins(ctx context.Context, titles []string, aggFuncTitle string) ([]entities.Coin, error) {

	if len(titles) == 0 {
		return []entities.Coin{}, nil
	}

	var query string
	switch aggFuncTitle {
	case "AVG", "MAX", "MIN", "SUM":
		query = fmt.Sprintf(`
		SELECT coin_name, %s(price) as price, MAX(created_at) as created_at
		FROM coins
		WHERE coin_name = ANY($1)
		GROUP BY coin_name
		`, aggFuncTitle)
	default:
		return nil, errors.New("unsupported aggregate function")
	}

	rows, err := s.db.Query(ctx, query, titles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute aggregate query")
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var coin entities.Coin
		if err := rows.Scan(&coin.CoinName, &coin.Price); err != nil {
			return nil, errors.Wrap(err, "failed to scan coin data")
		}
		coins = append(coins, coin)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error during rows iteration")
	}
	return coins, nil
}
