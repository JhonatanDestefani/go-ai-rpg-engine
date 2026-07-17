-- Align players table with the JSON-based player repository used by the RPG.

ALTER TABLE players
    ALTER COLUMN telegram_id DROP NOT NULL;

-- JSON blobs used by player_repository.go
ALTER TABLE players
    ADD COLUMN IF NOT EXISTS stats JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS skills JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS inventory JSONB NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE players
    ADD COLUMN IF NOT EXISTS equipped_weapon_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS equipped_armor_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS story_state JSONB NOT NULL DEFAULT '{}'::jsonb;

-- Legacy columnar stats from 001 are unused by Go; allow inserts that omit them.
ALTER TABLE players ALTER COLUMN max_hp SET DEFAULT 0;
ALTER TABLE players ALTER COLUMN hp SET DEFAULT 0;
ALTER TABLE players ALTER COLUMN max_mp SET DEFAULT 0;
ALTER TABLE players ALTER COLUMN mp SET DEFAULT 0;
ALTER TABLE players ALTER COLUMN attack SET DEFAULT 0;
ALTER TABLE players ALTER COLUMN defense SET DEFAULT 0;
ALTER TABLE players ALTER COLUMN magic SET DEFAULT 0;

ALTER TABLE players ALTER COLUMN max_hp DROP NOT NULL;
ALTER TABLE players ALTER COLUMN hp DROP NOT NULL;
ALTER TABLE players ALTER COLUMN max_mp DROP NOT NULL;
ALTER TABLE players ALTER COLUMN mp DROP NOT NULL;
ALTER TABLE players ALTER COLUMN attack DROP NOT NULL;
ALTER TABLE players ALTER COLUMN defense DROP NOT NULL;
ALTER TABLE players ALTER COLUMN magic DROP NOT NULL;
