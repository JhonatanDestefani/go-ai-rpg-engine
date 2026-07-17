CREATE TABLE player_inventory (
    id BIGSERIAL PRIMARY KEY,

    player_id BIGINT NOT NULL
        REFERENCES players(id)
        ON DELETE CASCADE,

    item_id BIGINT NOT NULL
        REFERENCES items(id)
        ON DELETE CASCADE,

    quantity INTEGER NOT NULL DEFAULT 1,
    is_equipped BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT player_inventory_quantity_check
        CHECK (quantity >= 1),

    CONSTRAINT player_inventory_unique_item
        UNIQUE (player_id, item_id)
);

CREATE INDEX idx_player_inventory_player
ON player_inventory(player_id);

CREATE INDEX idx_player_inventory_item
ON player_inventory(item_id);