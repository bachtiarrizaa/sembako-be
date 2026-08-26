CREATE TABLE product_discounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discount_id UUID NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
    product_id  UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_product_discount UNIQUE (discount_id, product_id)
);

CREATE INDEX idx_product_discounts_discount ON product_discounts(discount_id);
CREATE INDEX idx_product_discounts_product ON product_discounts(product_id);
