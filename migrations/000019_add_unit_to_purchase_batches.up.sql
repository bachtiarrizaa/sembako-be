ALTER TABLE purchase_batches
    ADD COLUMN unit_id UUID REFERENCES product_units(id),
    ADD COLUMN unit_price NUMERIC(14,2);