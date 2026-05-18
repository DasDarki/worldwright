# Worldwright

A self-hosted worldbuilding wiki for tabletop RPG settings — see `concept.md` for the full design.

## Quick start (Docker)

```bash
docker compose up --build
```

Then open <http://localhost:3000>. On first boot the backend runs migrations and seeds an admin account:

- Email: `admin@worldwright.local`
- Password: `worldwright` *(change it from the account page after first login)*

Set `WORLDWRIGHT_SESSION_SECRET` in a `.env` next to `docker-compose.yml` before exposing the service publicly.

For correct social-share previews (Discord, Slack, Twitter), also set in `.env`:

```
WORLDWRIGHT_PUBLIC_URL=https://worldwright.example.com
WORLDWRIGHT_SITE_NAME=Worldwright
```

`WORLDWRIGHT_PUBLIC_URL` is used both for OAuth callback URLs and to build absolute Open Graph URLs (`og:image`, `og:url`). Without it the SSR renderer falls back to the request's `Host` header — works on a single domain, fails for crawlers that don't follow redirects.

Fonts are downloaded at build time by `@nuxtjs/google-fonts` and bundled with the app; no runtime requests are made to Google's servers (GDPR-friendly).

### Progressive Web App

Worldwright ships as an installable PWA via `@vite-pwa/nuxt`:

- A web app manifest is generated at `/manifest.webmanifest` with name, theme color, icons (192/512 + maskable), and app shortcuts (Entries, Maps, Events, Calendars).
- A service worker precaches the app shell + fonts and runtime-caches uploaded asset images (StaleWhileRevalidate, 30-day TTL). API responses are **not** cached — auth-dependent data stays fresh.
- An install banner appears on supported browsers (Chrome/Edge desktop, Android, iOS-via-Safari "Add to Home Screen"). Update banners surface when a new build is detected (auto-update strategy with hourly poll).
- iOS-specific meta tags configure status-bar style, app title, and launch icons.

The PWA is **only enabled in production builds** (`devOptions.enabled: false`) — `npm run dev` serves regular pages without a service worker so hot-reload works cleanly. After `docker compose up --build`, open the app in Chrome and the install prompt should appear.

## Local development

### Backend (Go)

Requires Go 1.23+.

```bash
cd backend
go run ./cmd/server
```

The server listens on `:8080`. The SQLite file defaults to `./data/worldwright.db` (override via `WORLDWRIGHT_DATABASE_PATH`).

### Frontend (Nuxt 3)

Requires Node 20+ or [Bun](https://bun.sh) 1.1+. Bun is recommended — installs ~10× faster and runs the dev server natively:

```bash
cd frontend
bun install
bun run dev
```

…or with npm:

```bash
cd frontend
npm install --legacy-peer-deps
npm run dev
```

Vite serves at `http://localhost:3000` and proxies API calls to `http://localhost:8080`. The production Docker image uses Bun for both install + build (multi-stage) and runtime.

## Architecture

| Layer | Tech |
|---|---|
| API | Go · Fiber v2 · SQLite (`modernc.org/sqlite`) · goose migrations · bcrypt sessions |
| App | Nuxt 3 · SSR · Pinia · TipTap · `@nuxtjs/i18n` · `@nuxtjs/color-mode` · Tailwind + SCSS |
| Deploy | Two-container docker-compose; SQLite + assets on a named volume |

See `concept.md` for the full domain model, schema, and design language.

## Project layout

```
backend/    Go API, migrations, SQL queries
frontend/   Nuxt 3 SSR app
landing.html  Standalone marketing page (the cartographer's codex aesthetic source)
concept.md  Design + architecture spec
```
