BEGIN;

CREATE TABLE coins_new IF EXISTS(
    coin_name VARCHAR(50) PRIMARY KEY,
    price DECIMAL(15, 2) NOT NULL CHECK (price > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO coins_new (coin_name, price, created_at)
SELECT DISTINCT ON (coin_name) coin_name, price, created_at
FROM coins
ORDER BY coin_name, created_at DESC;

DROP TABLE coins;

ALTER TABLE coins_new RENAME TO coins;

COMMIT;