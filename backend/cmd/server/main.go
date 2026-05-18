package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"worldwright/backend/db"
	"worldwright/backend/internal/admin"
	"worldwright/backend/internal/auth"
	"worldwright/backend/internal/config"
	httpx "worldwright/backend/internal/http"
	"worldwright/backend/internal/store"
)

func main() {
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o755); err != nil {
		log.Error("create database dir", "err", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(cfg.AssetsDir, 0o755); err != nil {
		log.Error("create assets dir", "err", err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", cfg.DatabasePath)
	dbConn, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Error("open database", "err", err)
		os.Exit(1)
	}
	defer dbConn.Close()
	dbConn.SetMaxOpenConns(1)

	if err := dbConn.Ping(); err != nil {
		log.Error("ping database", "err", err)
		os.Exit(1)
	}

	goose.SetBaseFS(db.Migrations)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Error("goose dialect", "err", err)
		os.Exit(1)
	}
	if err := goose.Up(dbConn, "migrations"); err != nil {
		log.Error("goose up", "err", err)
		os.Exit(1)
	}

	st := store.New(dbConn)
	authSvc := auth.New(st, cfg.CookieSecure)
	oauthSvc := auth.NewOAuth(st, auth.OAuthSettings{
		PublicURL:           cfg.PublicURL,
		AllowSignup:         cfg.OAuthAllowSignup,
		GoogleClientID:      cfg.GoogleClientID,
		GoogleClientSecret:  cfg.GoogleClientSecret,
		DiscordClientID:     cfg.DiscordClientID,
		DiscordClientSecret: cfg.DiscordClientSecret,
	})

	if err := authSvc.EnsureAdmin(context.Background(), cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Error("ensure admin", "err", err)
		os.Exit(1)
	}
	adminSvc := admin.New(st, cfg.AssetsDir, log)
	log.Info("oauth providers", "enabled", oauthSvc.Enabled())

	app := fiber.New(fiber.Config{
		AppName:               "worldwright",
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		BodyLimit:             32 * 1024 * 1024,
		ErrorHandler:          httpx.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} ${method} ${path} ${latency}\n",
		TimeFormat: "15:04:05",
	}))
	if cfg.AllowOrigins != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     cfg.AllowOrigins,
			AllowMethods:     "GET,POST,PATCH,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Cookie",
			AllowCredentials: true,
		}))
	}

	httpx.RegisterRoutes(app, st, authSvc, oauthSvc, adminSvc, cfg.AssetsDir)

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.Listen(cfg.ListenAddr)
	}()

	log.Info("worldwright listening", "addr", cfg.ListenAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-listenErr:
		if err != nil {
			log.Error("server listen", "err", err)
			os.Exit(1)
		}
	case sig := <-quit:
		log.Info("shutdown", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Error("shutdown", "err", err)
		}
	}
}
