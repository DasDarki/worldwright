package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Store) CreateAsset(ctx context.Context, n NewAsset) (*Asset, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO assets (filename, path, mime, size, width, height, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.Filename, n.Path, n.Mime, n.Size, nullablePtrInt(n.Width), nullablePtrInt(n.Height), nowText())
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return s.AssetByID(ctx, id)
}

func (s *Store) AssetByID(ctx context.Context, id int64) (*Asset, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, filename, path, mime, size, width, height, created_at
		 FROM assets WHERE id = ?`, id)
	var a Asset
	var w, h sql.NullInt64
	var createdAt string
	if err := row.Scan(&a.ID, &a.Filename, &a.Path, &a.Mime, &a.Size, &w, &h, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if w.Valid {
		v := int(w.Int64)
		a.Width = &v
	}
	if h.Valid {
		v := int(h.Int64)
		a.Height = &v
	}
	a.CreatedAt = parseTime(createdAt)
	return &a, nil
}

func (s *Store) ListAssets(ctx context.Context, limit int) ([]Asset, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, filename, path, mime, size, width, height, created_at
		 FROM assets ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Asset, 0)
	for rows.Next() {
		var a Asset
		var w, h sql.NullInt64
		var createdAt string
		if err := rows.Scan(&a.ID, &a.Filename, &a.Path, &a.Mime, &a.Size, &w, &h, &createdAt); err != nil {
			return nil, err
		}
		if w.Valid {
			v := int(w.Int64)
			a.Width = &v
		}
		if h.Valid {
			v := int(h.Int64)
			a.Height = &v
		}
		a.CreatedAt = parseTime(createdAt)
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) DeleteAsset(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM assets WHERE id = ?`, id)
	return err
}
