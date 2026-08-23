CREATE TABLE transaction_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id      UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    product_unit_id     UUID NOT NULL REFERENCES product_units(id),
    qty                 NUMERIC(14,4) NOT NULL,
    unit_price          NUMERIC(14,2) NOT NULL,
    discount_applied    NUMERIC(14,2) NOT NULL DEFAULT 0,
    subtotal            NUMERIC(14,2) NOT NULL,
    total_cost          NUMERIC(14,2),
    margin              NUMERIC(14,2),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transaction_items_transaction ON transaction_items(transaction_id);
CREATE INDEX idx_transaction_items_product_unit ON transaction_items(product_unit_id);
