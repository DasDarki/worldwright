package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Store) UserByOAuth(ctx context.Context, provider, subject string) (*User, error) {
	var userID int64
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id FROM user_oauth_identities WHERE provider = ? AND subject = ?`,
		provider, subject).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.UserByID(ctx, userID)
}

func (s *Store) LinkOAuthIdentity(ctx context.Context, userID int64, provider, subject, email string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO user_oauth_identities (user_id, provider, subject, email, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		userID, provider, subject, nullable(email), nowText())
	return err
}

func (s *Store) ListOAuthIdentities(ctx context.Context, userID int64) ([]OAuthIdentity, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, user_id, provider, subject, COALESCE(email,''), created_at
		 FROM user_oauth_identities WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []OAuthIdentity
	for rows.Next() {
		var ident OAuthIdentity
		var createdAt string
		if err := rows.Scan(&ident.ID, &ident.UserID, &ident.Provider, &ident.Subject, &ident.Email, &createdAt); err != nil {
			return nil, err
		}
		ident.CreatedAt = parseTime(createdAt)
		out = append(out, ident)
	}
	return out, rows.Err()
}
