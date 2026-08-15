CREATE TABLE purchase_batches (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id          UUID NOT NULL REFERENCES products(id),
    supplier_id         UUID NOT NULL REFERENCES suppliers(id),
    initial_qty         NUMERIC(14,4) NOT NULL,   -- quantity awal dalam base unit
    remaining_qty       NUMERIC(14,4) NOT NULL,   -- kuantitas tersisa untuk FIFO
    purchase_price      NUMERIC(14,2) NOT NULL,   -- harga beli per base unit
    invoice_number      VARCHAR(100),
    purchase_date       DATE NOT NULL,
    created_by          UUID NOT NULL REFERENCES users(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_remaining_qty_valid CHECK (remaining_qty >= 0 AND remaining_qty <= initial_qty)
);

CREATE INDEX idx_purchase_batches_fifo ON purchase_batches(product_id, purchase_date, created_at) WHERE remaining_qty > 0;

CREATE TRIGGER trg_purchase_batches_updated_at BEFORE UPDATE ON purchase_batches
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
