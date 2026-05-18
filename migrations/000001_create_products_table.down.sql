-- 000001_create_products_table.down.sql

DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_products_category;
DROP TABLE IF EXISTS products;