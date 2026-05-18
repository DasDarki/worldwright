-- +goose Up

CREATE TABLE system_settings (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

INSERT INTO system_settings (key, value, updated_at) VALUES
  ('onboarding_completed', 'false', strftime('%Y-%m-%dT%H:%M:%fZ','now')),
  ('seed_data_present',    'true',  strftime('%Y-%m-%dT%H:%M:%fZ','now'));

ALTER TABLE entities      ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE relationships ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE tags          ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE calendars     ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE events        ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE maps          ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE timelines     ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE assets        ADD COLUMN is_seed INTEGER NOT NULL DEFAULT 0;

-- mark everything currently in the database as seed data so that the first
-- admin login can choose to keep or prune it. system entity_types and
-- relationship_types are infrastructure and are NOT flagged here.
UPDATE entities      SET is_seed = 1;
UPDATE relationships SET is_seed = 1;
UPDATE tags          SET is_seed = 1;
UPDATE calendars     SET is_seed = 1;
UPDATE events        SET is_seed = 1;
UPDATE maps          SET is_seed = 1;
UPDATE timelines     SET is_seed = 1;
UPDATE assets        SET is_seed = 1;

CREATE INDEX idx_entities_seed   ON entities(is_seed);
CREATE INDEX idx_tags_seed       ON tags(is_seed);
CREATE INDEX idx_calendars_seed  ON calendars(is_seed);

-- +goose Down

DROP TABLE IF EXISTS system_settings;
DROP INDEX IF EXISTS idx_entities_seed;
DROP INDEX IF EXISTS idx_tags_seed;
DROP INDEX IF EXISTS idx_calendars_seed;
-- SQLite doesn't support DROP COLUMN before 3.35; leave columns in place on rollback.
