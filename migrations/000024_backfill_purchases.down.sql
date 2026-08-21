-- Reverse of backfill: only removes headers that were auto-created by the
-- backfill migration (new purchases created after this migration are kept).
UPDATE purchase_batches pb
SET purchase_id = NULL
FROM purchases p
WHERE pb.purchase_id = p.id
  AND p.created_at = p.updated_at
  AND NOT EXISTS (
      SELECT 1 FROM purchases pp WHERE pp.id = p.id AND pp.updated_at > pp.created_at
  );

DELETE FROM purchases p
WHERE p.created_at = p.updated_at
  AND NOT EXISTS (SELECT 1 FROM purchase_batches pb WHERE pb.purchase_id = p.id);