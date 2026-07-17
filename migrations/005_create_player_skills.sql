CREATE TABLE player_skills (
    id BIGSERIAL PRIMARY KEY,

    player_id BIGINT NOT NULL
        REFERENCES players(id)
        ON DELETE CASCADE,

    skill_id BIGINT NOT NULL
        REFERENCES skills(id)
        ON DELETE CASCADE,

    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT player_skills_unique
        UNIQUE (player_id, skill_id)
);

CREATE INDEX idx_player_skills_player
ON player_skills(player_id);

CREATE INDEX idx_player_skills_skill
ON player_skills(skill_id);