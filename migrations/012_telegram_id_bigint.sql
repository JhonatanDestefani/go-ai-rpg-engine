-- Telegram chat/user IDs often exceed PostgreSQL integer (max ~2.1e9).
-- Force telegram_id to bigint even if an older DB created it as integer.
ALTER TABLE players
    ALTER COLUMN telegram_id TYPE BIGINT
    USING telegram_id::bigint;
