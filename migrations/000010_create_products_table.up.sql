CREATE TABLE products (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id                 UUID NOT NULL REFERENCES categories(id),
    name                        VARCHAR(150) NOT NULL,
    base_unit_id                UUID NOT NULL REFERENCES units(id),
    minimum_stock               NUMERIC(14,4),
    margin_threshold_percent    NUMERIC(5,2),
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_products_category ON products(category_id);

CREATE TRIGGER trg_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
