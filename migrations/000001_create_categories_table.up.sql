CREATE TABLE IF NOT EXISTS categories (
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
 
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories (deleted_at);
 
INSERT INTO categories (name, created_by, updated_by) VALUES
    ('Fruit',     'system', 'system'),
    ('Vegetable', 'system', 'system'),
    ('Dairy',     'system', 'system'),
    ('Meat',      'system', 'system')
ON CONFLICT (name) DO NOTHING;
 