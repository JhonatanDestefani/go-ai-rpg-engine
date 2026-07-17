CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,

    telegram_id BIGINT NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    class VARCHAR(30) NOT NULL,
    current_scene VARCHAR(100) NOT NULL DEFAULT 'start',

    level INTEGER NOT NULL DEFAULT 1,
    xp INTEGER NOT NULL DEFAULT 0,
    gold INTEGER NOT NULL DEFAULT 0,

    max_hp INTEGER NOT NULL,
    hp INTEGER NOT NULL,
    max_mp INTEGER NOT NULL,
    mp INTEGER NOT NULL,
    attack INTEGER NOT NULL,
    defense INTEGER NOT NULL,
    magic INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT players_level_check CHECK (level >= 1),
    CONSTRAINT players_xp_check CHECK (xp >= 0),
    CONSTRAINT players_gold_check CHECK (gold >= 0),
    CONSTRAINT players_hp_check CHECK (hp >= 0),
    CONSTRAINT players_mp_check CHECK (mp >= 0)
);

CREATE INDEX idx_players_telegram_id
ON players(telegram_id);