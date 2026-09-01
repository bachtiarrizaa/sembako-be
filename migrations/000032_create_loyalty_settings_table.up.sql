CREATE TABLE IF NOT EXISTS loyalty_settings (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  earning_rate     NUMERIC(14, 2) NOT NULL,
  redemption_rate  NUMERIC(14, 2) NOT NULL,
  minimum_redeem   INTEGER        NOT NULL,
  is_expiry_active BOOLEAN        NOT NULL DEFAULT FALSE,
  expiry_months    INTEGER        NOT NULL DEFAULT 12,
  created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
