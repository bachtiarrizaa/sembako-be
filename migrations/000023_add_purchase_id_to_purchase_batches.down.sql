DROP INDEX IF EXISTS idx_purchase_batches_purchase_id;
ALTER TABLE purchase_batches DROP COLUMN IF EXISTS purchase_id;