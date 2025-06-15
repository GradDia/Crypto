ALTER TABLE coins DROP CONSTRAINT coins_pkey;

CREATE INDEX idx_coins_coin_name ON coins(coin_name);