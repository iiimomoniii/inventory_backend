ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id BIGINT;

UPDATE products p
SET category_id = c.id
FROM categories c
WHERE c.name = p.category;

ALTER TABLE products
    ADD CONSTRAINT fk_products_category
    FOREIGN KEY (category_id)
    REFERENCES categories (id);

CREATE INDEX IF NOT EXISTS idx_products_category_id ON products (category_id);