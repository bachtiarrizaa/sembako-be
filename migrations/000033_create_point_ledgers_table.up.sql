CREATE TABLE IF NOT EXISTS point_ledgers (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id    UUID           NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  transaction_id UUID           REFERENCES transactions(id)       ON DELETE SET NULL,
  type           VARCHAR(20)    NOT NULL,
  points         INTEGER        NOT NULL,
  description    TEXT           NOT NULL DEFAULT '',
  expired_at     TIMESTAMPTZ,
  created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_point_ledgers_customer_id ON point_ledgers(customer_id);
CREATE INDEX IF NOT EXISTS idx_point_ledgers_type_expired_at ON point_ledgers(type, expired_at);
