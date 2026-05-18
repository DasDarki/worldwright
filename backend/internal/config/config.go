package config

import (
	"os"
	"strings"
)

type Config struct {
	ListenAddr    string
	DatabasePath  string
	AssetsDir     string
	SessionSecret string
	AllowOrigins  string
	AdminEmail    string
	AdminPassword string
	CookieSecure  bool

	PublicURL          string
	OAuthAllowSignup   bool
	GoogleClientID     string
	GoogleClientSecret string
	DiscordClientID    string
	DiscordClientSecret string
}

func Load() Config {
	publicURL := strings.TrimRight(envOr("WORLDWRIGHT_PUBLIC_URL", "http://localhost:3000"), "/")
	return Config{
		ListenAddr:    envOr("WORLDWRIGHT_LISTEN_ADDR", ":8080"),
		DatabasePath:  envOr("WORLDWRIGHT_DATABASE_PATH", "./data/worldwright.db"),
		AssetsDir:     envOr("WORLDWRIGHT_ASSETS_DIR", "./data/assets"),
		SessionSecret: envOr("WORLDWRIGHT_SESSION_SECRET", "dev-secret-change-me"),
		AllowOrigins:  envOr("WORLDWRIGHT_ALLOW_ORIGINS", ""),
		AdminEmail:    envOr("WORLDWRIGHT_ADMIN_EMAIL", "admin@worldwright.local"),
		AdminPassword: envOr("WORLDWRIGHT_ADMIN_PASSWORD", "worldwright"),
		CookieSecure:  os.Getenv("WORLDWRIGHT_COOKIE_SECURE") == "true",

		PublicURL:           publicURL,
		OAuthAllowSignup:    os.Getenv("WORLDWRIGHT_OAUTH_ALLOW_SIGNUP") == "true",
		GoogleClientID:      os.Getenv("WORLDWRIGHT_GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("WORLDWRIGHT_GOOGLE_CLIENT_SECRET"),
		DiscordClientID:     os.Getenv("WORLDWRIGHT_DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("WORLDWRIGHT_DISCORD_CLIENT_SECRET"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
