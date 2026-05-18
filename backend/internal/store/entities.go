package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type EntityFilter struct {
	EntityTypeID int64
	ParentID     *int64
	Visibility   []string
	Tag          string
}

func (s *Store) ListEntities(ctx context.Context, filter EntityFilter) ([]EntitySummary, error) {
	clauses := []string{"1=1"}
	args := []any{}
	if filter.EntityTypeID != 0 {
		clauses = append(clauses, "e.entity_type_id = ?")
		args = append(args, filter.EntityTypeID)
	}
	if filter.ParentID != nil {
		if *filter.ParentID == 0 {
			clauses = append(clauses, "e.parent_id IS NULL")
		} else {
			clauses = append(clauses, "e.parent_id = ?")
			args = append(args, *filter.ParentID)
		}
	}
	if len(filter.Visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(filter.Visibility)), ",")
		clauses = append(clauses, "e.visibility IN ("+placeholders+")")
		for _, v := range filter.Visibility {
			args = append(args, v)
		}
	}
	join := ""
	if filter.Tag != "" {
		join = `JOIN entity_tags et ON et.entity_id = e.id
		        JOIN tags t ON t.id = et.tag_id`
		clauses = append(clauses, "t.name = ?")
		args = append(args, strings.ToLower(filter.Tag))
	}
	query := fmt.Sprintf(
		`SELECT e.id, e.entity_type_id, e.title, e.slug, COALESCE(e.summary,''), e.parent_id, e.visibility
		 FROM entities e %s WHERE %s ORDER BY e.sort_order, e.title`,
		join, strings.Join(clauses, " AND "),
	)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]EntitySummary, 0)
	for rows.Next() {
		var e EntitySummary
		var parentID sql.NullInt64
		if err := rows.Scan(&e.ID, &e.EntityTypeID, &e.Title, &e.Slug, &e.Summary, &parentID, &e.Visibility); err != nil {
			return nil, err
		}
		if parentID.Valid {
			id := parentID.Int64
			e.ParentID = &id
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) EntityBySlug(ctx context.Context, slug string, visibility []string) (*Entity, error) {
	clause := ""
	args := []any{slug}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	row := s.db.QueryRowContext(ctx,
		`SELECT id, entity_type_id, title, slug, COALESCE(summary,''), body, parent_id, materialized_path, visibility, created_at, updated_at
		 FROM entities WHERE slug = ?`+clause, args...)
	return s.scanEntity(ctx, row)
}

func (s *Store) EntityByID(ctx context.Context, id int64, visibility []string) (*Entity, error) {
	clause := ""
	args := []any{id}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	row := s.db.QueryRowContext(ctx,
		`SELECT id, entity_type_id, title, slug, COALESCE(summary,''), body, parent_id, materialized_path, visibility, created_at, updated_at
		 FROM entities WHERE id = ?`+clause, args...)
	return s.scanEntity(ctx, row)
}

func (s *Store) scanEntity(ctx context.Context, row *sql.Row) (*Entity, error) {
	var e Entity
	var parentID sql.NullInt64
	var bodyText, createdAt, updatedAt string
	if err := row.Scan(&e.ID, &e.EntityTypeID, &e.Title, &e.Slug, &e.Summary, &bodyText, &parentID, &e.MaterializedPath, &e.Visibility, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if parentID.Valid {
		id := parentID.Int64
		e.ParentID = &id
	}
	e.Body = json.RawMessage(bodyText)
	e.CreatedAt = parseTime(createdAt)
	e.UpdatedAt = parseTime(updatedAt)
	e.Tags = make([]string, 0)
	e.FieldValues = map[string]string{}
	if err := s.loadEntityRelations(ctx, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *Store) loadEntityRelations(ctx context.Context, e *Entity) error {
	tagRows, err := s.db.QueryContext(ctx,
		`SELECT t.name FROM entity_tags et JOIN tags t ON t.id = et.tag_id
		 WHERE et.entity_id = ? ORDER BY t.name`, e.ID)
	if err != nil {
		return err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var name string
		if err := tagRows.Scan(&name); err != nil {
			return err
		}
		e.Tags = append(e.Tags, name)
	}

	fvRows, err := s.db.QueryContext(ctx,
		`SELECT fd.key, COALESCE(efv.value, '') FROM entity_field_values efv
		 JOIN field_definitions fd ON fd.id = efv.field_definition_id
		 WHERE efv.entity_id = ?`, e.ID)
	if err != nil {
		return err
	}
	defer fvRows.Close()
	for fvRows.Next() {
		var key, value string
		if err := fvRows.Scan(&key, &value); err != nil {
			return err
		}
		e.FieldValues[key] = value
	}
	return nil
}

func (s *Store) CreateEntity(ctx context.Context, n NewEntity) (*Entity, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	path, err := computeMaterializedPath(ctx, tx, n.ParentID)
	if err != nil {
		return nil, err
	}

	body := string(n.Body)
	if body == "" {
		body = "{}"
	}
	now := nowText()
	// New entities go to the end of their parent's child list.
	var maxOrder sql.NullInt64
	if err := tx.QueryRowContext(ctx,
		`SELECT MAX(sort_order) FROM entities WHERE COALESCE(parent_id, 0) = COALESCE(?, 0)`,
		nullableInt(n.ParentID)).Scan(&maxOrder); err != nil {
		return nil, err
	}
	sortOrder := int64(100)
	if maxOrder.Valid {
		sortOrder = maxOrder.Int64 + 100
	}
	res, err := tx.ExecContext(ctx,
		`INSERT INTO entities (entity_type_id, title, slug, summary, body, body_text, parent_id, materialized_path, visibility, sort_order, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.EntityTypeID, n.Title, n.Slug, nullable(n.Summary), body, n.BodyText,
		nullableInt(n.ParentID), path, n.Visibility, sortOrder, now, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	if err := syncEntityTags(ctx, tx, id, n.Tags); err != nil {
		return nil, err
	}
	if err := syncEntityFieldValues(ctx, tx, id, n.FieldValues, n.EntityTypeID); err != nil {
		return nil, err
	}
	if err := syncEntityLinks(ctx, tx, id, n.Wikilinks); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.EntityByID(ctx, id, nil)
}

func (s *Store) UpdateEntity(ctx context.Context, id int64, n NewEntity) (*Entity, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	path, err := computeMaterializedPath(ctx, tx, n.ParentID)
	if err != nil {
		return nil, err
	}
	body := string(n.Body)
	if body == "" {
		body = "{}"
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE entities SET entity_type_id = ?, title = ?, slug = ?, summary = ?, body = ?, body_text = ?,
		   parent_id = ?, materialized_path = ?, visibility = ?, updated_at = ?
		 WHERE id = ?`,
		n.EntityTypeID, n.Title, n.Slug, nullable(n.Summary), body, n.BodyText,
		nullableInt(n.ParentID), path, n.Visibility, nowText(), id); err != nil {
		return nil, err
	}

	if err := syncEntityTags(ctx, tx, id, n.Tags); err != nil {
		return nil, err
	}
	if err := syncEntityFieldValues(ctx, tx, id, n.FieldValues, n.EntityTypeID); err != nil {
		return nil, err
	}
	if err := syncEntityLinks(ctx, tx, id, n.Wikilinks); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.EntityByID(ctx, id, nil)
}

func (s *Store) DeleteEntity(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM entities WHERE id = ?`, id)
	return err
}

// ReorderEntry describes one row in a batch reorder request. ParentID is
// optional: when zero the entity stays under its current parent.
type ReorderEntry struct {
	ID        int64
	ParentID  *int64
	SortOrder int
}

// ReorderEntities updates sort_order (and optionally parent_id) for the given
// entities in a single transaction. Materialized paths are recomputed for any
// entity that changes parent, propagating to descendants.
func (s *Store) ReorderEntities(ctx context.Context, entries []ReorderEntry) error {
	if len(entries) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, e := range entries {
		if e.ParentID != nil {
			// reparent: recompute materialized_path, then propagate to descendants
			newPath, err := computeMaterializedPath(ctx, tx, e.ParentID)
			if err != nil {
				return err
			}
			var oldPath string
			if err := tx.QueryRowContext(ctx,
				`SELECT materialized_path FROM entities WHERE id = ?`, e.ID).Scan(&oldPath); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx,
				`UPDATE entities SET parent_id = ?, materialized_path = ?, sort_order = ?, updated_at = ?
				 WHERE id = ?`,
				nullableInt(e.ParentID), newPath, e.SortOrder, nowText(), e.ID); err != nil {
				return err
			}
			// shift descendants' materialized_path prefix
			descPrefixOld := oldPath + fmt.Sprintf("%d/", e.ID)
			descPrefixNew := newPath + fmt.Sprintf("%d/", e.ID)
			if _, err := tx.ExecContext(ctx,
				`UPDATE entities
				   SET materialized_path = ? || substr(materialized_path, length(?) + 1)
				 WHERE materialized_path LIKE ? || '%'`,
				descPrefixNew, descPrefixOld, descPrefixOld); err != nil {
				return err
			}
		} else {
			if _, err := tx.ExecContext(ctx,
				`UPDATE entities SET sort_order = ?, updated_at = ? WHERE id = ?`,
				e.SortOrder, nowText(), e.ID); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Store) Backlinks(ctx context.Context, entityID int64, visibility []string) ([]Backlink, error) {
	clause := ""
	args := []any{entityID}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND e.visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT e.id, e.title, e.slug, e.entity_type_id, COALESCE(e.summary, '')
		 FROM links l JOIN entities e ON e.id = l.source_entity_id
		 WHERE l.target_entity_id = ?`+clause+` ORDER BY e.title`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Backlink, 0)
	for rows.Next() {
		var bl Backlink
		if err := rows.Scan(&bl.SourceEntityID, &bl.Title, &bl.Slug, &bl.EntityTypeID, &bl.Summary); err != nil {
			return nil, err
		}
		out = append(out, bl)
	}
	return out, rows.Err()
}

func (s *Store) RecentEntities(ctx context.Context, visibility []string, limit int) ([]EntitySummary, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	clauses := []string{"1=1"}
	args := []any{}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clauses = append(clauses, "visibility IN ("+placeholders+")")
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	args = append(args, limit)
	query := `SELECT id, entity_type_id, title, slug, COALESCE(summary,''), parent_id, visibility
	          FROM entities WHERE ` + strings.Join(clauses, " AND ") + `
	          ORDER BY updated_at DESC LIMIT ?`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]EntitySummary, 0)
	for rows.Next() {
		var e EntitySummary
		var parentID sql.NullInt64
		if err := rows.Scan(&e.ID, &e.EntityTypeID, &e.Title, &e.Slug, &e.Summary, &parentID, &e.Visibility); err != nil {
			return nil, err
		}
		if parentID.Valid {
			id := parentID.Int64
			e.ParentID = &id
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) ListTags(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) ListTagsWithCounts(ctx context.Context, visibility []string) ([]TagWithCount, error) {
	clause := ""
	args := []any{}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " WHERE e.visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	query := `SELECT t.name, COUNT(DISTINCT e.id) AS c
	          FROM tags t
	          LEFT JOIN entity_tags et ON et.tag_id = t.id
	          LEFT JOIN entities e ON e.id = et.entity_id` + clause + `
	          GROUP BY t.id
	          ORDER BY c DESC, t.name`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]TagWithCount, 0)
	for rows.Next() {
		var t TagWithCount
		if err := rows.Scan(&t.Name, &t.Count); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func computeMaterializedPath(ctx context.Context, tx *sql.Tx, parentID *int64) (string, error) {
	if parentID == nil || *parentID == 0 {
		return "", nil
	}
	var parentPath string
	err := tx.QueryRowContext(ctx,
		`SELECT materialized_path FROM entities WHERE id = ?`, *parentID).Scan(&parentPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%d/", parentPath, *parentID), nil
}

func syncEntityTags(ctx context.Context, tx *sql.Tx, entityID int64, tags []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM entity_tags WHERE entity_id = ?`, entityID); err != nil {
		return err
	}
	for _, name := range tags {
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			continue
		}
		var tagID int64
		err := tx.QueryRowContext(ctx, `SELECT id FROM tags WHERE name = ?`, name).Scan(&tagID)
		if errors.Is(err, sql.ErrNoRows) {
			res, err := tx.ExecContext(ctx, `INSERT INTO tags (name) VALUES (?)`, name)
			if err != nil {
				return err
			}
			tagID, _ = res.LastInsertId()
		} else if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO entity_tags (entity_id, tag_id) VALUES (?, ?)`, entityID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func syncEntityFieldValues(ctx context.Context, tx *sql.Tx, entityID int64, values map[string]string, entityTypeID int64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM entity_field_values WHERE entity_id = ?`, entityID); err != nil {
		return err
	}
	for key, value := range values {
		if value == "" {
			continue
		}
		var fdID int64
		err := tx.QueryRowContext(ctx,
			`SELECT id FROM field_definitions WHERE entity_type_id = ? AND key = ?`,
			entityTypeID, key).Scan(&fdID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO entity_field_values (entity_id, field_definition_id, value) VALUES (?, ?, ?)`,
			entityID, fdID, value); err != nil {
			return err
		}
	}
	return nil
}

func syncEntityLinks(ctx context.Context, tx *sql.Tx, entityID int64, slugs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM links WHERE source_entity_id = ?`, entityID); err != nil {
		return err
	}
	for _, slug := range slugs {
		if slug == "" {
			continue
		}
		var targetID int64
		err := tx.QueryRowContext(ctx, `SELECT id FROM entities WHERE slug = ?`, slug).Scan(&targetID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return err
		}
		if targetID == entityID {
			continue
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO links (source_entity_id, target_entity_id, anchor) VALUES (?, ?, '')`,
			entityID, targetID); err != nil {
			return err
		}
	}
	return nil
}

func nullableInt(p *int64) any {
	if p == nil || *p == 0 {
		return nil
	}
	return *p
}
