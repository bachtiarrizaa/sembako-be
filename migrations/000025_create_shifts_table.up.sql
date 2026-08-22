CREATE TABLE shifts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cashier_id          UUID NOT NULL REFERENCES users(id),
    opening_balance     NUMERIC(14,2) NOT NULL,
    closing_balance     NUMERIC(14,2),                 -- actual cash counted at close
    system_balance      NUMERIC(14,2),                 -- opening_balance + total cash sales
    discrepancy         NUMERIC(14,2),                 -- closing_balance - system_balance
    discrepancy_note    TEXT,                          -- SHIFT-07: required if |discrepancy| > 1000
    status              VARCHAR(10) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'closed')),
    opened_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    closed_at           TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_discrepancy_note CHECK (
        status = 'open'
        OR discrepancy IS NULL
        OR abs(discrepancy) <= 1000
        OR discrepancy_note IS NOT NULL
    )
);

CREATE INDEX idx_shifts_cashier ON shifts(cashier_id, opened_at);

CREATE UNIQUE INDEX uq_one_open_shift_per_cashier
    ON shifts(cashier_id) WHERE status = 'open';

CREATE TRIGGER trg_shifts_updated_at BEFORE UPDATE ON shifts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
