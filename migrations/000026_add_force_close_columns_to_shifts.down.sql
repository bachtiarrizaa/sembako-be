ALTER TABLE shifts
  DROP COLUMN IF EXISTS force_close_reason,
  DROP COLUMN IF EXISTS force_closed_by;
