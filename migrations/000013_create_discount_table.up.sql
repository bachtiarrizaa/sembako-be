CREATE TABLE discounts (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                        VARCHAR(150) NOT NULL,
    type                        VARCHAR(10) NOT NULL CHECK (type IN ('percent', 'fixed')),
    value                       NUMERIC(14,2) NOT NULL,
    start_date                  DATE,
    end_date                    DATE,
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_discounts_active ON discounts(is_active) WHERE is_active = TRUE;

CREATE TRIGGER trg_discounts_updated_at BEFORE UPDATE ON discounts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
