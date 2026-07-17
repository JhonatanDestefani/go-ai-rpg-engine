CREATE TABLE skills (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,

    class VARCHAR(30),
    mana_cost INTEGER NOT NULL DEFAULT 0,
    base_damage INTEGER NOT NULL DEFAULT 0,
    required_level INTEGER NOT NULL DEFAULT 1,

    target_type VARCHAR(30) NOT NULL DEFAULT 'enemy',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT skills_mana_cost_check CHECK (mana_cost >= 0),
    CONSTRAINT skills_required_level_check CHECK (required_level >= 1)
);

CREATE INDEX idx_skills_code
ON skills(code);

CREATE INDEX idx_skills_class
ON skills(class);