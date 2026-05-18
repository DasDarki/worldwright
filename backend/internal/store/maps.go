package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type NewMap struct {
	Name           string
	AssetID        int64
	Width          int
	Height         int
	ParentEntityID *int64
}

func (s *Store) ListMaps(ctx context.Context) ([]MapAsset, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, asset_id, width, height, parent_entity_id
		 FROM maps ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]MapAsset, 0)
	for rows.Next() {
		var m MapAsset
		var parent sql.NullInt64
		if err := rows.Scan(&m.ID, &m.Name, &m.AssetID, &m.Width, &m.Height, &parent); err != nil {
			return nil, err
		}
		if parent.Valid {
			id := parent.Int64
			m.ParentEntityID = &id
		}
		m.Pins = []MapPin{}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) MapByID(ctx context.Context, id int64, visibility []string) (*MapAsset, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, asset_id, width, height, parent_entity_id
		 FROM maps WHERE id = ?`, id)
	var m MapAsset
	var parent sql.NullInt64
	if err := row.Scan(&m.ID, &m.Name, &m.AssetID, &m.Width, &m.Height, &parent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if parent.Valid {
		v := parent.Int64
		m.ParentEntityID = &v
	}

	clause := ""
	args := []any{id}
	if len(visibility) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(visibility)), ",")
		clause = " AND visibility IN (" + placeholders + ")"
		for _, v := range visibility {
			args = append(args, v)
		}
	}
	pinRows, err := s.db.QueryContext(ctx,
		`SELECT p.id, p.map_id, p.x, p.y, COALESCE(p.label,''), COALESCE(p.icon,''),
		    p.target_entity_id, COALESCE(e.slug, ''), p.target_map_id, p.visibility
		 FROM map_pins p
		 LEFT JOIN entities e ON e.id = p.target_entity_id
		 WHERE p.map_id = ?`+clause+` ORDER BY p.id`, args...)
	if err != nil {
		return nil, err
	}
	defer pinRows.Close()
	m.Pins = []MapPin{}
	for pinRows.Next() {
		var p MapPin
		var entID, mapID sql.NullInt64
		if err := pinRows.Scan(&p.ID, &p.MapID, &p.X, &p.Y, &p.Label, &p.Icon, &entID, &p.TargetEntitySlug, &mapID, &p.Visibility); err != nil {
			return nil, err
		}
		if entID.Valid {
			v := entID.Int64
			p.TargetEntityID = &v
		}
		if mapID.Valid {
			v := mapID.Int64
			p.TargetMapID = &v
		}
		m.Pins = append(m.Pins, p)
	}
	return &m, nil
}

func (s *Store) CreateMap(ctx context.Context, n NewMap) (*MapAsset, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO maps (name, asset_id, width, height, parent_entity_id)
		 VALUES (?, ?, ?, ?, ?)`,
		n.Name, n.AssetID, n.Width, n.Height, nullableInt(n.ParentEntityID))
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return s.MapByID(ctx, id, nil)
}

func (s *Store) UpdateMap(ctx context.Context, id int64, n NewMap) (*MapAsset, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE maps SET name = ?, width = ?, height = ?, parent_entity_id = ? WHERE id = ?`,
		n.Name, n.Width, n.Height, nullableInt(n.ParentEntityID), id)
	if err != nil {
		return nil, err
	}
	return s.MapByID(ctx, id, nil)
}

func (s *Store) DeleteMap(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM maps WHERE id = ?`, id)
	return err
}

func (s *Store) CreateMapPin(ctx context.Context, mapID int64, p NewMapPin) (*MapPin, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO map_pins (map_id, x, y, label, icon, target_entity_id, target_map_id, visibility)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		mapID, p.X, p.Y, nullable(p.Label), nullable(p.Icon),
		nullableInt(p.TargetEntityID), nullableInt(p.TargetMapID), p.Visibility)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return s.MapPinByID(ctx, id)
}

func (s *Store) UpdateMapPin(ctx context.Context, id int64, p NewMapPin) (*MapPin, error) {
	_, err := s.db.ExecContext(ctx,
		`UPDATE map_pins SET x = ?, y = ?, label = ?, icon = ?, target_entity_id = ?, target_map_id = ?, visibility = ?
		 WHERE id = ?`,
		p.X, p.Y, nullable(p.Label), nullable(p.Icon),
		nullableInt(p.TargetEntityID), nullableInt(p.TargetMapID), p.Visibility, id)
	if err != nil {
		return nil, err
	}
	return s.MapPinByID(ctx, id)
}

func (s *Store) DeleteMapPin(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM map_pins WHERE id = ?`, id)
	return err
}

func (s *Store) MapPinByID(ctx context.Context, id int64) (*MapPin, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, map_id, x, y, COALESCE(label,''), COALESCE(icon,''),
		    target_entity_id, target_map_id, visibility
		 FROM map_pins WHERE id = ?`, id)
	var p MapPin
	var entID, mapID sql.NullInt64
	if err := row.Scan(&p.ID, &p.MapID, &p.X, &p.Y, &p.Label, &p.Icon, &entID, &mapID, &p.Visibility); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if entID.Valid {
		v := entID.Int64
		p.TargetEntityID = &v
	}
	if mapID.Valid {
		v := mapID.Int64
		p.TargetMapID = &v
	}
	return &p, nil
}
