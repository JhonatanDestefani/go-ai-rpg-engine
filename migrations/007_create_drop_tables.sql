CREATE TABLE drop_tables (
    id BIGSERIAL PRIMARY KEY,

    monster_id BIGINT NOT NULL UNIQUE
        REFERENCES monsters(id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_drop_tables_monster
ON drop_tables(monster_id);