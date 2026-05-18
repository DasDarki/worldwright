-- +goose Up

CREATE TABLE calendar_moon_phases (
  id INTEGER PRIMARY KEY,
  moon_id INTEGER NOT NULL REFERENCES calendar_moons(id) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  name TEXT NOT NULL,
  cycle_position REAL NOT NULL,
  randomness REAL NOT NULL DEFAULT 0,
  icon TEXT
);

CREATE INDEX idx_calendar_moon_phases_moon ON calendar_moon_phases(moon_id);

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 1, 'New',             0.000, 0.020, '🌑' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 2, 'Waxing Crescent', 0.125, 0.050, '🌒' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 3, 'First Quarter',   0.250, 0.020, '🌓' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 4, 'Waxing Gibbous',  0.375, 0.050, '🌔' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 5, 'Full',            0.500, 0.020, '🌕' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 6, 'Waning Gibbous',  0.625, 0.050, '🌖' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 7, 'Last Quarter',    0.750, 0.020, '🌗' FROM calendar_moons m;

INSERT INTO calendar_moon_phases (moon_id, sort_order, name, cycle_position, randomness, icon)
SELECT m.id, 8, 'Waning Crescent', 0.875, 0.050, '🌘' FROM calendar_moons m;

-- +goose Down

DROP TABLE IF EXISTS calendar_moon_phases;
