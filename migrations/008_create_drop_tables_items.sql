CREATE TABLE drop_table_items (
    id BIGSERIAL PRIMARY KEY,

    drop_table_id BIGINT NOT NULL
        REFERENCES drop_tables(id)
        ON DELETE CASCADE,

    item_id BIGINT NOT NULL
        REFERENCES items(id)
        ON DELETE CASCADE,

    drop_chance NUMERIC(5,2) NOT NULL,
    min_quantity INTEGER NOT NULL DEFAULT 1,
    max_quantity INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT drop_chance_check
        CHECK (drop_chance >= 0 AND drop_chance <= 100),

    CONSTRAINT drop_quantity_check
        CHECK (
            min_quantity >= 1
            AND max_quantity >= min_quantity
        ),

    CONSTRAINT drop_table_item_unique
        UNIQUE (drop_table_id, item_id)
);

CREATE INDEX idx_drop_table_items_table
ON drop_table_items(drop_table_id);

CREATE INDEX idx_drop_table_items_item
ON drop_table_items(item_id);