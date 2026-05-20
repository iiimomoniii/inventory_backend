CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL      PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    category_id BIGINT         NOT NULL REFERENCES categories (id),
    price       DECIMAL(10, 2) NOT NULL CHECK (price > 0),
    stock       INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_by  VARCHAR(100),
    updated_by  VARCHAR(100),
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);
 
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products (category_id);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at  ON products (deleted_at);
 