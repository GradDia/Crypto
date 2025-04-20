package postgres

import (
	"Cryptoproject/internal/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"Cryptoproject/internal/cases"
)

type Storage struct {
	db *pgxpool.Pool
}

func newStorage(connectionString string) (*Storage, error) {
	con, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection pool")
	}

	return &Storage{db: con}, nil
}

func (s *Storage) Store(ctx context.Context, coins []entities.Coin) error {

}
