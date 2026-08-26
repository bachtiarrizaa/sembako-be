CREATE TABLE transaction_item_cost_allocations (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_item_id     UUID NOT NULL REFERENCES transaction_items(id) ON DELETE CASCADE,
    purchase_batch_id       UUID NOT NULL REFERENCES purchase_batches(id),
    qty_allocated           NUMERIC(14,4) NOT NULL,
    purchase_price_at_sale  NUMERIC(14,2) NOT NULL,
    cost_subtotal           NUMERIC(14,2) NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_cost_alloc_item ON transaction_item_cost_allocations(transaction_item_id);
CREATE INDEX idx_cost_alloc_batch ON transaction_item_cost_allocations(purchase_batch_id);
