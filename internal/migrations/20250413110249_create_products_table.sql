-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reception_id UUID NOT NULL REFERENCES receptions(id),
    type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reception_id ON products(reception_id);

-- +goose Down
DROP INDEX IF EXISTS idx_reception_id;
DROP TABLE IF EXISTS products;