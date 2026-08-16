-- 1. Cache agregat stok produk (1 row per product)
CREATE TABLE stocks (
    product_id      UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    qty_base_unit   NUMERIC(14,4) NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_stocks_updated_at BEFORE UPDATE ON stocks
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- 2. Log Kartu Stok (Append-only)
CREATE TABLE stock_mutations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id      UUID NOT NULL REFERENCES products(id),
    type            VARCHAR(10) NOT NULL CHECK (type IN ('in', 'out')),
    qty             NUMERIC(14,4) NOT NULL, -- selalu positif
    qty_before      NUMERIC(14,4) NOT NULL,
    qty_after       NUMERIC(14,4) NOT NULL,
    source          VARCHAR(30) NOT NULL CHECK (source IN ('purchase', 'sale', 'stock_count', 'damaged', 'lost', 'return')),
    reference_id    UUID, -- id dari purchase_batch atau stock_count
    note            TEXT,
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_stock_mutations_product ON stock_mutations(product_id, created_at);
CREATE INDEX idx_stock_mutations_reference ON stock_mutations(source, reference_id);
