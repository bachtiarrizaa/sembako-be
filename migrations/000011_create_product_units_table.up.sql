CREATE TABLE product_units (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id              UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    unit_id                 UUID NOT NULL REFERENCES units(id),
    conversion_to_base      NUMERIC(12,4) NOT NULL,
    selling_price           NUMERIC(14,2) NOT NULL,
    is_base_unit            BOOLEAN NOT NULL DEFAULT FALSE,
    is_active               BOOLEAN NOT NULL DEFAULT TRUE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, unit_id),
    CONSTRAINT chk_conversion_positive CHECK (conversion_to_base > 0)
);

CREATE INDEX idx_product_units_product ON product_units(product_id);
CREATE UNIQUE INDEX uq_one_base_unit_per_product
    ON product_units(product_id) WHERE is_base_unit = TRUE;

CREATE TRIGGER trg_product_units_updated_at BEFORE UPDATE ON product_units
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
