CREATE TABLE skill_effects (
    id BIGSERIAL PRIMARY KEY,

    skill_id BIGINT NOT NULL
        REFERENCES skills(id)
        ON DELETE CASCADE,

    effect_type VARCHAR(50) NOT NULL,
    duration INTEGER NOT NULL DEFAULT 0,
    power INTEGER NOT NULL DEFAULT 0,

    is_debuff BOOLEAN NOT NULL DEFAULT false,
    is_buff BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT skill_effects_duration_check
        CHECK (duration >= 0),

    CONSTRAINT skill_effects_buff_debuff_check
        CHECK (NOT (is_buff = true AND is_debuff = true))
);

CREATE INDEX idx_skill_effects_skill
ON skill_effects(skill_id);

CREATE INDEX idx_skill_effects_type
ON skill_effects(effect_type);