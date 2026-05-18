package store

import (
	"context"
	"database/sql"
)

// PruneContent deletes content from the database. When seedOnly is true,
// only rows flagged with is_seed = 1 are removed; user-created data stays
// untouched. When seedOnly is false, ALL content is wiped — this is used by
// import flows that replace the whole codex.
//
// Infrastructure (entity_types, relationship_types, field_definitions, users,
// sessions, oauth identities, system_settings) is preserved in both modes.
//
// Returns the disk-relative asset paths that were removed so the caller can
// unlink the files from the assets directory.
func (s *Store) PruneContent(ctx context.Context, seedOnly bool) ([]string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	assetPaths, err := selectAssetPaths(ctx, tx, seedOnly)
	if err != nil {
		return nil, err
	}

	if seedOnly {
		err = pruneSeedOnly(ctx, tx)
	} else {
		err = pruneAll(ctx, tx)
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return assetPaths, nil
}

func pruneAll(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`DELETE FROM timeline_events`,
		`DELETE FROM timelines`,
		`DELETE FROM event_entities`,
		`DELETE FROM events`,
		`DELETE FROM map_pins`,
		`DELETE FROM maps`,
		`DELETE FROM relationships`,
		`DELETE FROM links`,
		`DELETE FROM entity_tags`,
		`DELETE FROM entity_field_values`,
		`DELETE FROM entities`,
		`DELETE FROM calendar_moon_phases`,
		`DELETE FROM calendar_moons`,
		`DELETE FROM calendar_months`,
		`DELETE FROM calendar_weekdays`,
		`DELETE FROM calendar_eras`,
		`DELETE FROM calendars`,
		`DELETE FROM tags`,
		`DELETE FROM assets`,
	}
	for _, q := range stmts {
		if _, err := tx.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func pruneSeedOnly(ctx context.Context, tx *sql.Tx) error {
	stmts := []string{
		`DELETE FROM timeline_events WHERE timeline_id IN (SELECT id FROM timelines WHERE is_seed = 1)
		                                OR event_id    IN (SELECT id FROM events    WHERE is_seed = 1)`,
		`DELETE FROM timelines WHERE is_seed = 1`,

		`DELETE FROM event_entities WHERE event_id IN (SELECT id FROM events WHERE is_seed = 1)`,
		`DELETE FROM events WHERE is_seed = 1`,

		`DELETE FROM map_pins WHERE map_id IN (SELECT id FROM maps WHERE is_seed = 1)`,
		`DELETE FROM maps WHERE is_seed = 1`,

		`DELETE FROM relationships WHERE is_seed = 1
		                              OR from_entity_id IN (SELECT id FROM entities WHERE is_seed = 1)
		                              OR to_entity_id   IN (SELECT id FROM entities WHERE is_seed = 1)`,
		`DELETE FROM links WHERE source_entity_id IN (SELECT id FROM entities WHERE is_seed = 1)
		                      OR target_entity_id IN (SELECT id FROM entities WHERE is_seed = 1)`,
		`DELETE FROM entity_tags         WHERE entity_id IN (SELECT id FROM entities WHERE is_seed = 1)`,
		`DELETE FROM entity_field_values WHERE entity_id IN (SELECT id FROM entities WHERE is_seed = 1)`,
		`DELETE FROM entities WHERE is_seed = 1`,

		// calendar children cascade via FK; only drop the seed calendar if nothing else references it
		`DELETE FROM calendars WHERE is_seed = 1
		                        AND id NOT IN (SELECT calendar_id FROM events)`,

		// tags + assets: leave them if user-created content still references them
		`DELETE FROM tags   WHERE is_seed = 1 AND id NOT IN (SELECT tag_id   FROM entity_tags)`,
		`DELETE FROM assets WHERE is_seed = 1 AND id NOT IN (SELECT asset_id FROM maps)`,
	}
	for _, q := range stmts {
		if _, err := tx.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func selectAssetPaths(ctx context.Context, tx *sql.Tx, seedOnly bool) ([]string, error) {
	q := `SELECT path FROM assets`
	if seedOnly {
		q += ` WHERE is_seed = 1 AND id NOT IN (SELECT asset_id FROM maps WHERE is_seed = 0)`
	}
	rows, err := tx.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
