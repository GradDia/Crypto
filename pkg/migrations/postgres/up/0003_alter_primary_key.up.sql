BEGIN;

-- Создаем новую таблицу с автоинкрементным ID и без unique constraint
CREATE TABLE IF NOT EXISTS coins (
                                     id SERIAL PRIMARY KEY,
                                     coin_name VARCHAR(50) NOT NULL,
    price DECIMAL(15, 2) NOT NULL CHECK (price > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Создаем индекс для поиска (но не unique!)
CREATE INDEX IF NOT EXISTS idx_coins_coin_name ON coins(coin_name);

COMMIT;