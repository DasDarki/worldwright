package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

func (s *Store) ListEntityTypes(ctx context.Context) ([]EntityType, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, key, name_en, name_de, COALESCE(icon,''), COALESCE(color,''), is_system
		 FROM entity_types ORDER BY name_en`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	types := make([]EntityType, 0)
	indexByID := map[int64]int{}
	for rows.Next() {
		var t EntityType
		var isSystem int
		if err := rows.Scan(&t.ID, &t.Key, &t.NameEn, &t.NameDe, &t.Icon, &t.Color, &isSystem); err != nil {
			return nil, err
		}
		t.IsSystem = isSystem != 0
		t.Fields = make([]FieldDefinition, 0)
		indexByID[t.ID] = len(types)
		types = append(types, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	fieldRows, err := s.db.QueryContext(ctx,
		`SELECT id, entity_type_id, key, label_en, label_de, data_type, COALESCE(config,''), sort_order, is_required
		 FROM field_definitions ORDER BY entity_type_id, sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer fieldRows.Close()
	for fieldRows.Next() {
		var f FieldDefinition
		var configText string
		var required int
		if err := fieldRows.Scan(&f.ID, &f.EntityTypeID, &f.Key, &f.LabelEn, &f.LabelDe, &f.DataType, &configText, &f.SortOrder, &required); err != nil {
			return nil, err
		}
		f.IsRequired = required != 0
		if configText != "" {
			f.Config = json.RawMessage(configText)
		}
		if idx, ok := indexByID[f.EntityTypeID]; ok {
			types[idx].Fields = append(types[idx].Fields, f)
		}
	}
	return types, fieldRows.Err()
}

func (s *Store) EntityTypeByID(ctx context.Context, id int64) (*EntityType, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, key, name_en, name_de, COALESCE(icon,''), COALESCE(color,''), is_system
		 FROM entity_types WHERE id = ?`, id)
	var t EntityType
	var isSystem int
	if err := row.Scan(&t.ID, &t.Key, &t.NameEn, &t.NameDe, &t.Icon, &t.Color, &isSystem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	t.IsSystem = isSystem != 0
	return &t, nil
}
