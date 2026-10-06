DROP INDEX IF EXISTS sessions_hostname_idx;
DROP INDEX IF EXISTS sessions_family_id_idx;

ALTER TABLE sessions
DROP CONSTRAINT IF EXISTS sessions_hostname_valid,
DROP COLUMN IF EXISTS hostname,
DROP COLUMN IF EXISTS family_id;