CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,

    item_type VARCHAR(30) NOT NULL,
    rarity VARCHAR(30) NOT NULL,

    value INTEGER NOT NULL DEFAULT 0,
    stackable BOOLEAN NOT NULL DEFAULT false,
    max_stack INTEGER NOT NULL DEFAULT 1,

    attack_bonus INTEGER NOT NULL DEFAULT 0,
    defense_bonus INTEGER NOT NULL DEFAULT 0,
    magic_bonus INTEGER NOT NULL DEFAULT 0,
    heal_amount INTEGER NOT NULL DEFAULT 0,
    mana_amount INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT items_value_check CHECK (value >= 0),
    CONSTRAINT items_max_stack_check CHECK (max_stack >= 1)
);

CREATE INDEX idx_items_code
ON items(code);

CREATE INDEX idx_items_type
ON items(item_type);

CREATE INDEX idx_items_rarity
ON items(rarity);