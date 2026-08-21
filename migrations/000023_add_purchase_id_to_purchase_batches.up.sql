ALTER TABLE purchase_batches
    ADD COLUMN purchase_id UUID REFERENCES purchases(id);

CREATE INDEX idx_purchase_batches_purchase_id ON purchase_batches(purchase_id);