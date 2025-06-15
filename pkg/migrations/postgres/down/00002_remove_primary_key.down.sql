DROP INDEX idx_coins_coin_name;

ALTER TABLE coins ADD PRIMARY KEY (coin_name);