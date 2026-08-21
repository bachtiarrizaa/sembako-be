-- Backfill: group existing purchase batches into purchase headers.
-- Batches sharing (invoice_number, supplier_id, purchase_date, created_by)
-- are considered a single purchase (invoice). Batches with no invoice number
-- are grouped by (supplier_id, purchase_date, created_by).
INSERT INTO purchases (invoice_number, supplier_id, purchase_date, total_amount, created_by, created_at, updated_at)
SELECT
    invoice_number,
    supplier_id,
    purchase_date,
    SUM(initial_qty * purchase_price) AS total_amount,
    created_by,
    MIN(created_at) AS created_at,
    MIN(created_at) AS updated_at
FROM purchase_batches
WHERE purchase_id IS NULL
GROUP BY invoice_number, supplier_id, purchase_date, created_by;

UPDATE purchase_batches pb
SET purchase_id = p.id
FROM purchases p
WHERE pb.purchase_id IS NULL
  AND pb.invoice_number IS NOT DISTINCT FROM p.invoice_number
  AND pb.supplier_id = p.supplier_id
  AND pb.purchase_date = p.purchase_date
  AND pb.created_by = p.created_by;