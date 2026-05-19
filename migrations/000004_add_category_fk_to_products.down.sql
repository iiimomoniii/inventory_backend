DROP INDEX IF EXISTS idx_products_category_id;
ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_products_category;
ALTER TABLE products DROP COLUMN IF EXISTS category_id;