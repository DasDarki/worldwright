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
	GoogleClientID     string
	GoogleClientSecret string
	DiscordClientID    string
	DiscordClientSecret string

	OnyxURL          string
	OnyxClientID     string
	OnyxClientSecret string
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
		GoogleClientID:      os.Getenv("WORLDWRIGHT_GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("WORLDWRIGHT_GOOGLE_CLIENT_SECRET"),
		DiscordClientID:     os.Getenv("WORLDWRIGHT_DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("WORLDWRIGHT_DISCORD_CLIENT_SECRET"),

		OnyxURL:          strings.TrimRight(os.Getenv("WORLDWRIGHT_ONYX_URL"), "/"),
		OnyxClientID:     os.Getenv("WORLDWRIGHT_ONYX_CLIENT_ID"),
		OnyxClientSecret: os.Getenv("WORLDWRIGHT_ONYX_CLIENT_SECRET"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
