package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"Cryptoproject/internal/entities"
)

type Storage struct {
	db *pgxpool.Pool
}

func newStorage(connectionString string) (*Storage, error) {
	con, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
	}
}
