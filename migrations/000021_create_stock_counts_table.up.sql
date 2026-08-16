-- Transaksi Pengajuan Opname Fisik
CREATE TABLE stock_counts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id      UUID NOT NULL REFERENCES products(id),
    count_date      DATE NOT NULL DEFAULT current_date,
    system_qty      NUMERIC(14,4) NOT NULL, -- snapshot stok sistem saat submit
    physical_qty    NUMERIC(14,4) NOT NULL, -- input fisik dari kasir/admin
    discrepancy     NUMERIC(14,4) NOT NULL, -- physical_qty - system_qty
    note            TEXT,
    status          VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    submitted_by    UUID NOT NULL REFERENCES users(id),
    approved_by     UUID REFERENCES users(id),
    submitted_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_stock_counts_status ON stock_counts(status);
CREATE TRIGGER trg_stock_counts_updated_at BEFORE UPDATE ON stock_counts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
