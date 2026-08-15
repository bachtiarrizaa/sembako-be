ALTER TABLE purchase_batches
    DROP COLUMN IF EXISTS unit_price,
    DROP COLUMN IF EXISTS unit_id;