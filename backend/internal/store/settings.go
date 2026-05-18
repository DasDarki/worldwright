package store

import (
	"context"
	"database/sql"
	"errors"
)

type SystemSettings struct {
	OnboardingCompleted bool `json:"onboarding_completed"`
	SeedDataPresent     bool `json:"seed_data_present"`
}

func (s *Store) GetSystemSettings(ctx context.Context) (SystemSettings, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT key, value FROM system_settings`)
	if err != nil {
		return SystemSettings{}, err
	}
	defer rows.Close()
	out := SystemSettings{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return SystemSettings{}, err
		}
		switch k {
		case "onboarding_completed":
			out.OnboardingCompleted = v == "true"
		case "seed_data_present":
			out.SeedDataPresent = v == "true"
		}
	}
	return out, rows.Err()
}

func (s *Store) SetSystemSetting(ctx context.Context, key, value string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO system_settings (key, value, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
		key, value, nowText())
	return err
}

func (s *Store) GetSystemSetting(ctx context.Context, key string) (string, error) {
	var v string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = ?`, key).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return v, err
}
