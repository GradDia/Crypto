-- Чистый SQL без директив Goose
CREATE TABLE IF NOT EXISTS coins (
                                     coin_name VARCHAR(50) PRIMARY KEY,
    price DECIMAL(15, 2) NOT NULL CHECK (price > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_coins_created_at ON coins(created_at);

-- Триггер для updated_at
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_timestamp
    BEFORE UPDATE ON coins
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();