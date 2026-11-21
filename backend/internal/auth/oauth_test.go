package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"worldwright/backend/db"
	"worldwright/backend/internal/store"
)

func openDatabase(t *testing.T) *sql.DB {
	t.Helper()
	path := filepath.Join(t.TempDir(), "worldwright.db")
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", path)
	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	t.Cleanup(func() { database.Close() })
	goose.SetBaseFS(db.Migrations)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}
	return database
}

func TestOptionalEmailMigrationKeepsUsersAndReferences(t *testing.T) {
	database := openDatabase(t)
	if err := goose.UpTo(database, "migrations", 6); err != nil {
		t.Fatal(err)
	}
	statements := []string{
		`INSERT INTO users (id, email, password_hash, role, locale, display_name, created_at) VALUES (1, 'admin@worldwright.local', 'hash', 'admin', 'en', 'Admin', '2026-01-01T00:00:00Z')`,
		`INSERT INTO users (id, email, password_hash, role, locale, display_name, created_at) VALUES (2, 'player@worldwright.local', 'hash', 'player', 'de', 'Player', '2026-01-02T00:00:00Z')`,
		`INSERT INTO sessions (token, user_id, expires_at, created_at) VALUES ('token', 2, '2099-01-01T00:00:00Z', '2026-01-02T00:00:00Z')`,
		`INSERT INTO user_oauth_identities (user_id, provider, subject, email, created_at) VALUES (2, 'discord', '42', 'player@worldwright.local', '2026-01-02T00:00:00Z')`,
	}
	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			t.Fatal(err)
		}
	}

	if err := goose.Up(database, "migrations"); err != nil {
		t.Fatal(err)
	}

	st := store.New(database)
	ctx := context.Background()
	player, err := st.UserByID(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if player.Email != "player@worldwright.local" || player.Role != "player" || player.Locale != "de" || player.DisplayName != "Player" || !player.HasPassword {
		t.Fatalf("player changed during migration: %+v", player)
	}
	if linked, err := st.UserByOAuth(ctx, "discord", "42"); err != nil || linked.ID != 2 {
		t.Fatalf("oauth identity lost: %v %+v", err, linked)
	}
	var sessions int
	if err := database.QueryRow(`SELECT COUNT(*) FROM sessions`).Scan(&sessions); err != nil || sessions != 1 {
		t.Fatalf("sessions lost: %v %d", err, sessions)
	}

	for range 2 {
		created, err := st.CreateUser(ctx, store.NewUser{Role: "player", Locale: "en", DisplayName: "Onyx"})
		if err != nil {
			t.Fatalf("user without email: %v", err)
		}
		if created.Email != "" {
			t.Fatalf("expected empty email, got %q", created.Email)
		}
	}
	if _, err := st.CreateUser(ctx, store.NewUser{Email: "player@worldwright.local", Role: "player", Locale: "en"}); err == nil {
		t.Fatal("duplicate email accepted")
	}
	if _, err := database.Exec(`INSERT INTO sessions (token, user_id, expires_at, created_at) VALUES ('orphan', 999, '2099-01-01T00:00:00Z', '2026-01-02T00:00:00Z')`); err == nil {
		t.Fatal("foreign key to users is not enforced after the rebuild")
	}
	if _, err := database.Exec(`DELETE FROM users WHERE id = 2`); err != nil {
		t.Fatal(err)
	}
	if err := database.QueryRow(`SELECT COUNT(*) FROM sessions`).Scan(&sessions); err != nil || sessions != 0 {
		t.Fatalf("cascade delete broken: %v %d", err, sessions)
	}
}

func TestOAuthProvidersNoLongerCreateAccounts(t *testing.T) {
	database := openDatabase(t)
	if err := goose.Up(database, "migrations"); err != nil {
		t.Fatal(err)
	}
	st := store.New(database)
	ctx := context.Background()
	existing, err := st.CreateUser(ctx, store.NewUser{Email: "known@worldwright.local", Role: "player", Locale: "en"})
	if err != nil {
		t.Fatal(err)
	}
	oauth := NewOAuth(st, OAuthSettings{PublicURL: "http://localhost:3000"})

	if _, err := oauth.ResolveUser(ctx, "discord", &OAuthUser{Subject: "new", Email: "stranger@example.com", Name: "Stranger"}); !errors.Is(err, ErrSignupDisabled) {
		t.Fatalf("expected signup to be disabled, got %v", err)
	}
	count, err := st.CountUsers(ctx)
	if err != nil || count != 1 {
		t.Fatalf("unexpected user count %d: %v", count, err)
	}

	resolved, err := oauth.ResolveUser(ctx, "google", &OAuthUser{Subject: "google-1", Email: "known@worldwright.local"})
	if err != nil || resolved.ID != existing.ID {
		t.Fatalf("existing account was not linked by email: %v", err)
	}
}
