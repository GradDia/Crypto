package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) Store(ctx context.Context, coins []entities.Coin) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	query := `
		INSERT INTO coins (coin_name, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (coin_name) 
		DO UPDATE SET 
			price = EXCLUDED.price,
			updated_at = EXCLUDED.updated_at
	`

	for _, coin := range coins {
		_, err = tx.ExecContext(ctx, query,
			coin.CoinName,
			coin.Price,
			coin.CreatedAt,
			time.Now(),
			)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "failed to insert coin")
		}
	}
	return tx.Commit()

}

func (p *PostgresStorage) GetCoinsList(ctx context.Context) ([]string, error) {
	query := `SELECT coin_name FROM coins ORDER BY coin_name`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query coins list")
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, errors.Wrap(err, "failed to scan coins list")
		}
		coins = append(coins, name)
	}
	return coins, nil
}

func
