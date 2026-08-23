ALTER TABLE shifts
  ADD COLUMN force_close_reason TEXT,
  ADD COLUMN force_closed_by UUID REFERENCES users(id);
