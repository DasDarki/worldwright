package store

import (
	"context"
	"database/sql"
	"errors"
)

func (s *Store) CountUsers(ctx context.Context) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

func (s *Store) CreateUser(ctx context.Context, n NewUser) (*User, error) {
	now := nowText()
	var pwHash sql.NullString
	if n.PasswordHash != "" {
		pwHash = sql.NullString{String: n.PasswordHash, Valid: true}
	}
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, role, locale, display_name, avatar_url, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.Email, pwHash, n.Role, n.Locale, nullable(n.DisplayName), nullable(n.AvatarURL), now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return s.UserByID(ctx, id)
}

func (s *Store) UserByID(ctx context.Context, id int64) (*User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, email, COALESCE(password_hash,''), role, locale, COALESCE(display_name,''), COALESCE(avatar_url,''), created_at
		 FROM users WHERE id = ?`, id)
	return scanUser(row)
}

func (s *Store) UserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, email, COALESCE(password_hash,''), role, locale, COALESCE(display_name,''), COALESCE(avatar_url,''), created_at
		 FROM users WHERE email = ?`, email)
	return scanUser(row)
}

func (s *Store) UpdateUserPassword(ctx context.Context, userID int64, passwordHash string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, userID)
	return err
}

func (s *Store) UpdateUserProfile(ctx context.Context, userID int64, displayName, avatarURL, locale string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET display_name = ?, avatar_url = ?, locale = ? WHERE id = ?`,
		nullable(displayName), nullable(avatarURL), locale, userID)
	return err
}

func scanUser(row *sql.Row) (*User, error) {
	var u User
	var createdAt string
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.Locale, &u.DisplayName, &u.AvatarURL, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.HasPassword = u.PasswordHash != ""
	u.CreatedAt = parseTime(createdAt)
	return &u, nil
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}
