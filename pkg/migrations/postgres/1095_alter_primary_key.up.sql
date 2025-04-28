BEGIN;

CREATE TABLE coins_new (
    coin_name VARCHAR(50) NOT NULL,
    price DECIMAL(15, 2) NOT NULL CHECK (price > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
)

INSERT INTO coins_new (coin_name, price, created_at)
SELECT coin_name, price, created_at FROM coins;

DROP TABLE coins;

ALTER TABLE coins_new RENAME TO coins;

COMMIT;