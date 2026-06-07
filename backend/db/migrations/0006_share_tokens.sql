-- +goose Up

CREATE TABLE share_tokens (
  token TEXT PRIMARY KEY,
  entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL,
  created_by INTEGER REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_share_tokens_expires ON share_tokens(expires_at);
CREATE INDEX idx_share_tokens_entity  ON share_tokens(entity_id);

-- +goose Down

DROP INDEX IF EXISTS idx_share_tokens_expires;
DROP INDEX IF EXISTS idx_share_tokens_entity;
DROP TABLE IF EXISTS share_tokens;
