package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

func (s *Store) ListRelationshipTypes(ctx context.Context) ([]RelationshipType, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, key, label_en, label_de,
		    COALESCE(inverse_label_en, ''), COALESCE(inverse_label_de, ''),
		    is_symmetric, category
		 FROM relationship_types
		 ORDER BY category, label_en`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]RelationshipType, 0)
	for rows.Next() {
		var rt RelationshipType
		var sym int
		if err := rows.Scan(&rt.ID, &rt.Key, &rt.LabelEn, &rt.LabelDe,
			&rt.InverseLabelEn, &rt.InverseLabelDe, &sym, &rt.Category); err != nil {
			return nil, err
		}
		rt.IsSymmetric = sym != 0
		out = append(out, rt)
	}
	return out, rows.Err()
}

func (s *Store) RelationshipTypeByID(ctx context.Context, id int64) (*RelationshipType, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, key, label_en, label_de,
		    COALESCE(inverse_label_en, ''), COALESCE(inverse_label_de, ''),
		    is_symmetric, category
		 FROM relationship_types WHERE id = ?`, id)
	var rt RelationshipType
	var sym int
	if err := row.Scan(&rt.ID, &rt.Key, &rt.LabelEn, &rt.LabelDe,
		&rt.InverseLabelEn, &rt.InverseLabelDe, &sym, &rt.Category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	rt.IsSymmetric = sym != 0
	return &rt, nil
}

func (s *Store) ListRelationshipsForEntity(ctx context.Context, entityID int64, visibility []string) ([]RelationshipEdge, error) {
	vClause := ""
	args := []any{entityID, entityID}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		vClause = " AND other.visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	query := `
		SELECT r.id,
		       CASE WHEN r.from_entity_id = ? THEN 'out' ELSE 'in' END AS direction,
		       other.id, other.entity_type_id, other.title, other.slug,
		       COALESCE(other.summary, ''), other.parent_id, other.visibility,
		       rt.id, rt.key, rt.label_en, rt.label_de,
		       COALESCE(rt.inverse_label_en, ''), COALESCE(rt.inverse_label_de, ''),
		       rt.is_symmetric, rt.category,
		       COALESCE(r.description, '')
		FROM relationships r
		JOIN relationship_types rt ON rt.id = r.relationship_type_id
		JOIN entities other
		  ON other.id = CASE WHEN r.from_entity_id = ? THEN r.to_entity_id ELSE r.from_entity_id END
		WHERE (r.from_entity_id = ? OR r.to_entity_id = ?)` + vClause + `
		ORDER BY rt.category, rt.label_en, other.title`
	args = append([]any{entityID, entityID}, args...)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]RelationshipEdge, 0)
	for rows.Next() {
		var e RelationshipEdge
		var parentID sql.NullInt64
		var sym int
		if err := rows.Scan(&e.ID, &e.Direction,
			&e.Other.ID, &e.Other.EntityTypeID, &e.Other.Title, &e.Other.Slug,
			&e.Other.Summary, &parentID, &e.Other.Visibility,
			&e.Type.ID, &e.Type.Key, &e.Type.LabelEn, &e.Type.LabelDe,
			&e.Type.InverseLabelEn, &e.Type.InverseLabelDe, &sym, &e.Type.Category,
			&e.Description); err != nil {
			return nil, err
		}
		e.Type.IsSymmetric = sym != 0
		if parentID.Valid {
			id := parentID.Int64
			e.Other.ParentID = &id
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

type NewRelationship struct {
	FromEntityID       int64
	ToEntityID         int64
	RelationshipTypeID int64
	Description        string
}

func (s *Store) CreateRelationship(ctx context.Context, n NewRelationship) (*Relationship, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO relationships (from_entity_id, to_entity_id, relationship_type_id, description)
		 VALUES (?, ?, ?, ?)`,
		n.FromEntityID, n.ToEntityID, n.RelationshipTypeID, nullable(n.Description))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &Relationship{
		ID:                 id,
		FromEntityID:       n.FromEntityID,
		ToEntityID:         n.ToEntityID,
		RelationshipTypeID: n.RelationshipTypeID,
		Description:        n.Description,
	}, nil
}

func (s *Store) DeleteRelationship(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM relationships WHERE id = ?`, id)
	return err
}

func (s *Store) RelationshipExists(ctx context.Context, fromID, toID, typeID int64) (bool, error) {
	var n int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM relationships
		 WHERE relationship_type_id = ?
		   AND ((from_entity_id = ? AND to_entity_id = ?) OR (from_entity_id = ? AND to_entity_id = ?))`,
		typeID, fromID, toID, toID, fromID).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
