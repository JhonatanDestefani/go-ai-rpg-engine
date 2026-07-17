ALTER TABLE players
ADD COLUMN equipped_weapon_id TEXT NOT NULL DEFAULT '',
ADD COLUMN equipped_armor_id TEXT NOT NULL DEFAULT '',
ADD COLUMN story_state JSONB NOT NULL DEFAULT '{}'::jsonb;