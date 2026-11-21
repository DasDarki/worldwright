-- +goose NO TRANSACTION
-- +goose Up
PRAGMA foreign_keys = OFF;

BEGIN;

CREATE TABLE users_rebuilt (
  id INTEGER PRIMARY KEY,
  email TEXT UNIQUE,
  password_hash TEXT,
  role TEXT NOT NULL DEFAULT 'player',
  locale TEXT NOT NULL DEFAULT 'en',
  display_name TEXT,
  avatar_url TEXT,
  created_at TEXT NOT NULL
);

INSERT INTO users_rebuilt (id, email, password_hash, role, locale, display_name, avatar_url, created_at)
SELECT id, email, password_hash, role, locale, display_name, avatar_url, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_rebuilt RENAME TO users;

COMMIT;

PRAGMA foreign_keys = ON;

-- +goose Down
PRAGMA foreign_keys = OFF;

BEGIN;

CREATE TABLE users_rebuilt (
  id INTEGER PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT,
  role TEXT NOT NULL DEFAULT 'player',
  locale TEXT NOT NULL DEFAULT 'en',
  display_name TEXT,
  avatar_url TEXT,
  created_at TEXT NOT NULL
);

INSERT INTO users_rebuilt (id, email, password_hash, role, locale, display_name, avatar_url, created_at)
SELECT id, COALESCE(email, 'user-' || id || '@worldwright.invalid'), password_hash, role, locale, display_name, avatar_url, created_at FROM users;

DROP TABLE users;

ALTER TABLE users_rebuilt RENAME TO users;

COMMIT;

PRAGMA foreign_keys = ON;
