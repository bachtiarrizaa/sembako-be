DROP TRIGGER IF EXISTS trg_shifts_updated_at ON shifts;
DROP INDEX IF EXISTS uq_one_open_shift_per_cashier;
DROP INDEX IF EXISTS idx_shifts_cashier;
DROP TABLE IF EXISTS shifts;
