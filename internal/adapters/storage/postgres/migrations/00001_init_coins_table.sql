-- +goose Up
-- SQL-запросы для наката миграции
CREATE TABLE coins (
                       coin_name VARCHAR(50) PRIMARY KEY,
                       price DECIMAL(15, 2) NOT NULL CHECK (price > 0),
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_coins_created_at ON coins(created_at);

-- +goose Down
-- SQL-запросы для отката миграции
DROP TABLE IF EXISTS coins;
DROP FUNCTION IF EXISTS update_timestamp CASCADE;