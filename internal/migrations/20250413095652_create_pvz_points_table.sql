-- +goose Up
CREATE TABLE IF NOT EXISTS pvz_points (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    city TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_created_at ON pvz_points(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_created_at;
DROP TABLE IF EXISTS pickup_points;