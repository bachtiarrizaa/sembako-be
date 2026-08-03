DROP TRIGGER IF EXISTS trg_product_units_updated_at ON product_units;
DROP INDEX IF EXISTS uq_one_base_unit_per_product;
DROP INDEX IF EXISTS idx_product_units_product;
DROP TABLE IF EXISTS product_units;
