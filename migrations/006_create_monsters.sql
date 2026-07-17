CREATE TABLE monsters (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,

    level INTEGER NOT NULL DEFAULT 1,

    max_hp INTEGER NOT NULL,
    max_mp INTEGER NOT NULL DEFAULT 0,
    attack INTEGER NOT NULL,
    defense INTEGER NOT NULL,
    magic INTEGER NOT NULL DEFAULT 0,

    xp_reward INTEGER NOT NULL DEFAULT 0,
    gold_min INTEGER NOT NULL DEFAULT 0,
    gold_max INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT monsters_level_check CHECK (level >= 1),
    CONSTRAINT monsters_max_hp_check CHECK (max_hp >= 1),
    CONSTRAINT monsters_xp_check CHECK (xp_reward >= 0),
    CONSTRAINT monsters_gold_check CHECK (
        gold_min >= 0
        AND gold_max >= gold_min
    )
);

CREATE INDEX idx_monsters_code
ON monsters(code);

CREATE INDEX idx_monsters_level
ON monsters(level);