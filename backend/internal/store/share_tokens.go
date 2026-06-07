package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"
)

type ShareToken struct {
	Token     string    `json:"token"`
	EntityID  int64     `json:"entity_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *int64    `json:"created_by,omitempty"`
}

// CreateShareToken mints a fresh, opaque token that lets anyone with the
// link view the entity until expiresAt regardless of the entity's
// visibility setting.
func (s *Store) CreateShareToken(ctx context.Context, entityID int64, createdBy int64, ttl time.Duration) (*ShareToken, error) {
	token, err := generateOpaqueToken(24)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	expires := now.Add(ttl)
	var byArg any
	if createdBy > 0 {
		byArg = createdBy
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO share_tokens (token, entity_id, expires_at, created_at, created_by)
		 VALUES (?, ?, ?, ?, ?)`,
		token, entityID, expires.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano), byArg)
	if err != nil {
		return nil, err
	}
	return &ShareToken{
		Token: token, EntityID: entityID,
		ExpiresAt: expires, CreatedAt: now,
	}, nil
}

// ResolveShareToken looks up the token and returns the linked entity ID if
// it is still valid. Expired tokens are treated as not found and lazily
// reaped.
func (s *Store) ResolveShareToken(ctx context.Context, token string) (int64, error) {
	if token == "" {
		return 0, ErrNotFound
	}
	var entityID int64
	var expiresStr string
	err := s.db.QueryRowContext(ctx,
		`SELECT entity_id, expires_at FROM share_tokens WHERE token = ?`, token,
	).Scan(&entityID, &expiresStr)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	exp, _ := time.Parse(time.RFC3339Nano, expiresStr)
	if exp.IsZero() || time.Now().UTC().After(exp) {
		// best-effort cleanup; ignore error
		_, _ = s.db.ExecContext(ctx, `DELETE FROM share_tokens WHERE token = ?`, token)
		return 0, ErrNotFound
	}
	return entityID, nil
}

// PurgeExpiredShareTokens removes any token past its expiry. Safe to call
// periodically; returns the number of rows removed.
func (s *Store) PurgeExpiredShareTokens(ctx context.Context) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`DELETE FROM share_tokens WHERE expires_at < ?`,
		time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// ListShareTokensFor returns all non-expired tokens currently bound to a
// given entity. Used by the admin UI to show / revoke active links.
func (s *Store) ListShareTokensFor(ctx context.Context, entityID int64) ([]ShareToken, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT token, entity_id, expires_at, created_at, created_by
		 FROM share_tokens
		 WHERE entity_id = ? AND expires_at > ?
		 ORDER BY created_at DESC`,
		entityID, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]ShareToken, 0)
	for rows.Next() {
		var t ShareToken
		var expiresStr, createdStr string
		var by sql.NullInt64
		if err := rows.Scan(&t.Token, &t.EntityID, &expiresStr, &createdStr, &by); err != nil {
			return nil, err
		}
		t.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expiresStr)
		t.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdStr)
		if by.Valid {
			v := by.Int64
			t.CreatedBy = &v
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) DeleteShareToken(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM share_tokens WHERE token = ?`, token)
	return err
}

func generateOpaqueToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
