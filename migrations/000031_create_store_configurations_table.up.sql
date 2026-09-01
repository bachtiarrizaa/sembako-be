CREATE TABLE IF NOT EXISTS store_configurations (
  id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  store_name                  VARCHAR(100) NOT NULL,
  store_address               TEXT,
  store_phone                 VARCHAR(20),
  receipt_header_text         TEXT,
  receipt_footer_text         TEXT,
  receipt_show_cashier_name   BOOLEAN NOT NULL DEFAULT TRUE,
  receipt_show_customer_name  BOOLEAN NOT NULL DEFAULT TRUE,
  shift_discrepancy_tolerance NUMERIC(14, 2) NOT NULL DEFAULT 1000,
  created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
