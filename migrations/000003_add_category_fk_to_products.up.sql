CREATE TABLE IF NOT EXISTS categories (
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories (deleted_at);

INSERT INTO categories (name) VALUES
    ('Fruit'),
    ('Vegetable'),
    ('Dairy'),
    ('Meat')
ON CONFLICT (name) DO NOTHING;