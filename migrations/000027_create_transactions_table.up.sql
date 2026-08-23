CREATE TABLE transactions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receipt_number          VARCHAR(30) NOT NULL UNIQUE,
    cashier_id              UUID NOT NULL REFERENCES users(id),
    shift_id                UUID NOT NULL REFERENCES shifts(id),
    customer_id             UUID REFERENCES customers(id),
    payment_method          VARCHAR(10) NOT NULL CHECK (payment_method IN ('cash', 'qris', 'transfer')),
    subtotal                NUMERIC(14,2) NOT NULL,
    total_discount          NUMERIC(14,2) NOT NULL DEFAULT 0,
    points_used             INTEGER NOT NULL DEFAULT 0,
    points_discount_value   NUMERIC(14,2) NOT NULL DEFAULT 0,
    points_earned           INTEGER NOT NULL DEFAULT 0,
    total                   NUMERIC(14,2) NOT NULL,
    cash_received           NUMERIC(14,2),
    change_given            NUMERIC(14,2),
    manual_paid_confirmation BOOLEAN,
    status                  VARCHAR(10) NOT NULL DEFAULT 'completed' CHECK (status IN ('completed', 'void')),
    void_reason             TEXT,
    voided_by               UUID REFERENCES users(id),
    voided_at               TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_customer_required_for_points CHECK (
        customer_id IS NOT NULL OR (points_used = 0 AND points_earned = 0)
    )
);

CREATE INDEX idx_transactions_shift ON transactions(shift_id);
CREATE INDEX idx_transactions_cashier_created ON transactions(cashier_id, created_at);
CREATE INDEX idx_transactions_customer ON transactions(customer_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

CREATE TRIGGER trg_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
