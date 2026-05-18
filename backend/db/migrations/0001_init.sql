-- +goose Up

CREATE TABLE entity_types (
  id INTEGER PRIMARY KEY,
  key TEXT UNIQUE NOT NULL,
  name_en TEXT NOT NULL,
  name_de TEXT NOT NULL,
  icon TEXT,
  color TEXT,
  is_system INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE field_definitions (
  id INTEGER PRIMARY KEY,
  entity_type_id INTEGER NOT NULL REFERENCES entity_types(id) ON DELETE CASCADE,
  key TEXT NOT NULL,
  label_en TEXT NOT NULL,
  label_de TEXT NOT NULL,
  data_type TEXT NOT NULL,
  config TEXT,
  sort_order INTEGER NOT NULL DEFAULT 0,
  is_required INTEGER NOT NULL DEFAULT 0,
  UNIQUE (entity_type_id, key)
);

CREATE TABLE entities (
  id INTEGER PRIMARY KEY,
  entity_type_id INTEGER NOT NULL REFERENCES entity_types(id),
  title TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  summary TEXT,
  body TEXT NOT NULL DEFAULT '{}',
  body_text TEXT NOT NULL DEFAULT '',
  parent_id INTEGER REFERENCES entities(id) ON DELETE SET NULL,
  materialized_path TEXT NOT NULL DEFAULT '',
  visibility TEXT NOT NULL DEFAULT 'secret',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE INDEX idx_entities_type ON entities(entity_type_id);
CREATE INDEX idx_entities_parent ON entities(parent_id);
CREATE INDEX idx_entities_path ON entities(materialized_path);
CREATE INDEX idx_entities_visibility ON entities(visibility);

CREATE TABLE entity_field_values (
  entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  field_definition_id INTEGER NOT NULL REFERENCES field_definitions(id) ON DELETE CASCADE,
  value TEXT,
  PRIMARY KEY (entity_id, field_definition_id)
);

CREATE TABLE calendars (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  epoch_name TEXT,
  current_year INTEGER NOT NULL DEFAULT 0,
  current_month_index INTEGER NOT NULL DEFAULT 0,
  current_day INTEGER NOT NULL DEFAULT 1,
  leap_rule TEXT
);

CREATE TABLE calendar_months (
  id INTEGER PRIMARY KEY,
  calendar_id INTEGER NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL,
  name TEXT NOT NULL,
  days INTEGER NOT NULL
);

CREATE TABLE calendar_weekdays (
  id INTEGER PRIMARY KEY,
  calendar_id INTEGER NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL,
  name TEXT NOT NULL
);

CREATE TABLE calendar_eras (
  id INTEGER PRIMARY KEY,
  calendar_id INTEGER NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  start_year INTEGER NOT NULL,
  suffix TEXT
);

CREATE TABLE calendar_moons (
  id INTEGER PRIMARY KEY,
  calendar_id INTEGER NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  cycle_days REAL NOT NULL,
  offset_days REAL NOT NULL DEFAULT 0
);

CREATE TABLE events (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  body TEXT NOT NULL DEFAULT '{}',
  calendar_id INTEGER NOT NULL REFERENCES calendars(id),
  era_id INTEGER REFERENCES calendar_eras(id) ON DELETE SET NULL,
  year INTEGER NOT NULL,
  month_index INTEGER NOT NULL,
  day INTEGER NOT NULL,
  end_year INTEGER,
  end_month_index INTEGER,
  end_day INTEGER,
  importance INTEGER NOT NULL DEFAULT 0,
  visibility TEXT NOT NULL DEFAULT 'secret'
);

CREATE TABLE event_entities (
  event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  role TEXT NOT NULL DEFAULT 'participant',
  PRIMARY KEY (event_id, entity_id, role)
);

CREATE TABLE relationship_types (
  id INTEGER PRIMARY KEY,
  key TEXT UNIQUE NOT NULL,
  label_en TEXT NOT NULL,
  label_de TEXT NOT NULL,
  inverse_label_en TEXT,
  inverse_label_de TEXT,
  is_symmetric INTEGER NOT NULL DEFAULT 0,
  category TEXT NOT NULL DEFAULT 'generic'
);

CREATE TABLE relationships (
  id INTEGER PRIMARY KEY,
  from_entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  to_entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  relationship_type_id INTEGER NOT NULL REFERENCES relationship_types(id),
  description TEXT,
  since_event_id INTEGER REFERENCES events(id) ON DELETE SET NULL
);

CREATE INDEX idx_relationships_from ON relationships(from_entity_id);
CREATE INDEX idx_relationships_to ON relationships(to_entity_id);

CREATE TABLE tags (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE entity_tags (
  entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (entity_id, tag_id)
);

CREATE TABLE links (
  source_entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  target_entity_id INTEGER NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
  anchor TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (source_entity_id, target_entity_id, anchor)
);

CREATE INDEX idx_links_target ON links(target_entity_id);

CREATE TABLE timelines (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT
);

CREATE TABLE timeline_events (
  timeline_id INTEGER NOT NULL REFERENCES timelines(id) ON DELETE CASCADE,
  event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  PRIMARY KEY (timeline_id, event_id)
);

CREATE TABLE assets (
  id INTEGER PRIMARY KEY,
  filename TEXT NOT NULL,
  path TEXT NOT NULL,
  mime TEXT NOT NULL,
  size INTEGER NOT NULL,
  width INTEGER,
  height INTEGER,
  created_at TEXT NOT NULL
);

CREATE TABLE maps (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  asset_id INTEGER NOT NULL REFERENCES assets(id),
  width INTEGER NOT NULL,
  height INTEGER NOT NULL,
  parent_entity_id INTEGER REFERENCES entities(id) ON DELETE SET NULL
);

CREATE TABLE map_pins (
  id INTEGER PRIMARY KEY,
  map_id INTEGER NOT NULL REFERENCES maps(id) ON DELETE CASCADE,
  x REAL NOT NULL,
  y REAL NOT NULL,
  label TEXT,
  icon TEXT,
  target_entity_id INTEGER REFERENCES entities(id) ON DELETE SET NULL,
  target_map_id INTEGER REFERENCES maps(id) ON DELETE SET NULL,
  visibility TEXT NOT NULL DEFAULT 'secret'
);

CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT,
  role TEXT NOT NULL DEFAULT 'player',
  locale TEXT NOT NULL DEFAULT 'en',
  display_name TEXT,
  avatar_url TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE user_oauth_identities (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider TEXT NOT NULL,
  subject TEXT NOT NULL,
  email TEXT,
  created_at TEXT NOT NULL,
  UNIQUE (provider, subject)
);

CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE INDEX idx_sessions_user ON sessions(user_id);

CREATE VIRTUAL TABLE entities_fts USING fts5(
  title, summary, body_text,
  content='entities', content_rowid='id'
);

-- +goose StatementBegin
CREATE TRIGGER entities_fts_after_insert AFTER INSERT ON entities BEGIN
  INSERT INTO entities_fts(rowid, title, summary, body_text)
  VALUES (new.id, new.title, COALESCE(new.summary, ''), new.body_text);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER entities_fts_after_update AFTER UPDATE ON entities BEGIN
  UPDATE entities_fts
    SET title = new.title, summary = COALESCE(new.summary, ''), body_text = new.body_text
    WHERE rowid = new.id;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER entities_fts_after_delete AFTER DELETE ON entities BEGIN
  DELETE FROM entities_fts WHERE rowid = old.id;
END;
-- +goose StatementEnd

-- +goose Down

DROP TABLE IF EXISTS entities_fts;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS user_oauth_identities;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS map_pins;
DROP TABLE IF EXISTS maps;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS timeline_events;
DROP TABLE IF EXISTS timelines;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS entity_tags;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS relationships;
DROP TABLE IF EXISTS relationship_types;
DROP TABLE IF EXISTS event_entities;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS calendar_moons;
DROP TABLE IF EXISTS calendar_eras;
DROP TABLE IF EXISTS calendar_weekdays;
DROP TABLE IF EXISTS calendar_months;
DROP TABLE IF EXISTS calendars;
DROP TABLE IF EXISTS entity_field_values;
DROP TABLE IF EXISTS entities;
DROP TABLE IF EXISTS field_definitions;
DROP TABLE IF EXISTS entity_types;
