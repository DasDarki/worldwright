-- +goose Up

ALTER TABLE entities ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;

-- Initialize existing rows so they have a stable order matching the previous
-- alphabetical ordering. Multiplying by 100 leaves room for future inserts.
UPDATE entities SET sort_order = (
  SELECT (COUNT(*) - 1) * 100
  FROM entities AS e2
  WHERE COALESCE(e2.parent_id, 0) = COALESCE(entities.parent_id, 0)
    AND e2.title <= entities.title
);

CREATE INDEX idx_entities_sort ON entities(parent_id, sort_order);

-- +goose Down

DROP INDEX IF EXISTS idx_entities_sort;
-- SQLite < 3.35 cannot drop columns; leave sort_order in place on rollback.
