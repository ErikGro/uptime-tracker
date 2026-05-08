# Roadmap

Ordered milestones for delivering the Principled Uptime Tracker as specified in `CLAUDE.md`. Each milestone is a coherent shippable slice; check items off as they land.

## Configuration boundary (decided up front)

- **`.env` (static, set at deploy time):** admin username, admin password, listen address, SQLite DB path, optional log level.
- **SQLite `settings` table (runtime, edited via the web UI):** poll interval, request timeout, failure threshold, webhook URL, webhook enabled flag, retention days.

The dashboard never lets you change auth credentials — those are infrastructure concerns and stay in `.env`. Everything else the operator might tune at 2am is one click away.

---

## Milestone 0 — Repo hygiene

- [x] Vendor `pico.min.css` into `internal/web/static/css/`.
- [x] Delete `pico-main/`.
- [x] Vendor `htmx.min.js` (htmx 4.0.0-beta2 — see [`docs/htmx-4-migration.md`](./docs/htmx-4-migration.md)) into `internal/web/static/js/`.
- [x] Vendor `alpine.min.js` (Alpine 3.15.11) into `internal/web/static/js/`.
- [x] Fill `.gitignore` (`*.db`, `.env`, `tmp/`, build artefacts).
- [x] Replace empty `README.md` with orientation + config reference (will keep expanding as features land).
- [x] Run `go mod tidy`; direct deps now include `gorm.io/gorm`, `gorm.io/driver/sqlite`, `a-h/templ`.
- [ ] **Deployment caveat:** `mattn/go-sqlite3` requires CGO, which conflicts with CLAUDE.md §6's `CGO_ENABLED=0` target. Swap to `modernc.org/sqlite` (pure Go) at M9, or accept CGO in the build pipeline.

**Done when:** repo contains only files we own or have explicitly vendored, with versions recorded in `README.md`.

## Milestone 1 — Configuration & settings store

- [x] `internal/config`: load `.env` with `os.Getenv` + a tiny loader (no `godotenv` dependency).
- [x] `settings` table schema: `(key TEXT PRIMARY KEY, value TEXT)`.
- [x] Typed accessors per key with defaults:
  - `poll_interval_seconds` (default `300`)
  - `request_timeout_seconds` (default `10`)
  - `failure_threshold` (default `3`)
  - `webhook_url` (nullable string)
  - `webhook_enabled` (default `false`)
  - `retention_days` (default `30`)
- [x] Seed defaults on first boot (idempotent via `ON CONFLICT DO NOTHING`).
- [x] `.env.example` checked in as a copy-this-and-edit template.

**Done when:** the app reads `.env` for static config and the `settings` table for runtime config, with sane defaults if either is missing. ✅

## Milestone 2 — Persistence layer

- [x] `internal/store`: GORM models —
  - `URL { ID, Label, URL, CurrentStatus, ConsecutiveFailures, LastCheckedAt, CreatedAt, UpdatedAt }`
  - `Check { ID, URLID, CheckedAt, StatusCode, LatencyMs, OK, Error }`
  - `Setting { Key, Value }`
- [x] `AutoMigrate` on startup.
- [x] CRUD repository functions in `urls.go`, `checks.go`, `settings.go`; no GORM types leak above this package.
- [x] Unit tests in `store_test.go` cover settings, URL CRUD, check append/prune, and status updates.

**Done when:** unit tests against SQLite can create/read/update/delete URLs and append checks. ✅

## Milestone 3 — HTTP skeleton

- [x] `internal/web`: HTTP mux with Basic Auth middleware reading creds from `config` (constant-time compare).
- [x] Embedded static FS via `//go:embed all:static`, served at `/static/` with `http.FileServerFS`.
- [x] Base `templ` layout linking `pico.min.css`, `htmx.min.js`, `alpine.min.js` from `/static/`.
- [x] `/healthz` returns 200 OK (unauthenticated, for deploy health checks).
- [x] `/static/*` unauthenticated; everything else gated by Basic Auth.

**Done when:** `go run .` boots, prompts for Basic Auth, and serves a Pico-styled page. ✅

## Milestone 4 — URL CRUD via HTMX

- [ ] `GET /` → dashboard with URL list (templ).
- [ ] `POST /urls` → returns the new row fragment.
- [ ] `GET /urls/{id}/edit` → returns edit form fragment.
- [ ] `PUT /urls/{id}` → returns updated row fragment.
- [ ] `DELETE /urls/{id}` → returns empty fragment + `HX-Trigger`.
- [ ] Server-side validation; surface errors via `hx-target` swap.

**Done when:** an operator can add, edit, and remove URLs without a full page reload.

## Milestone 5 — Scheduler / worker

- [ ] `internal/scheduler`: single goroutine, single `time.Ticker` at the global interval.
- [ ] On each tick: fan out checks across all URLs with `errgroup` + bounded concurrency.
- [ ] Per-URL: HTTP GET with timeout from settings, write `Check` row, update `URL.consecutive_failures`.
- [ ] State machine: `UP → DOWN` after `failure_threshold` consecutive failures; `DOWN → UP` on first success. Emit a `StateTransition` event on changes.
- [ ] Hot-reload: when settings change, the scheduler picks up the new interval/timeout/threshold without restart (e.g., a `chan struct{}` reload signal).

**Done when:** stopping the network on a monitored host flips it to `DOWN` after exactly three failed checks, and restoring it flips back to `UP` on the next success.

## Milestone 6 — Notifications (webhook)

- [ ] `internal/notifier`: webhook sender with retry/backoff (`cenkalti/backoff` already in `go.sum`).
- [ ] Subscribe to `StateTransition` events; POST JSON `{ url, label, status, latency_ms, transitioned_at }` to `webhook_url` when `webhook_enabled` is true.
- [ ] "Send test webhook" button on the Settings page.

**Done when:** a state transition fires exactly one POST to the configured URL, with retry on 5xx/network errors.

## Milestone 7 — Settings page

- [ ] `GET /settings` → form pre-filled from the `settings` table.
- [ ] `PUT /settings` → persist, signal scheduler reload, return success toast fragment.
- [ ] Read-only display of admin username so operators remember it lives in `.env`.

**Done when:** every runtime setting listed in Milestone 1 is editable from the UI and takes effect without restart.

## Milestone 8 — Dashboard polish

- [ ] Auto-refresh status list with `hx-trigger="every 10s"` (server-rendered, no client state).
- [ ] Per-URL detail page: recent checks table + inline SVG latency sparkline (server-rendered from the last N rows — no charting library).
- [ ] Empty-state and error fragments.
- [ ] Retention worker: nightly prune of `checks` older than `retention_days`.

**Done when:** the dashboard updates itself, drilling into a URL shows its history, and the DB doesn't grow unboundedly.

## Milestone 9 — Deployment

- [ ] `Makefile` targets: `build` (`CGO_ENABLED=0 go build`), `run`, `tidy`, `templ` (regenerate templates).
- [ ] Sample `Caddyfile` for HTTPS termination + reverse proxy.
- [ ] Sample `systemd` unit.
- [ ] Backup note: which file is the DB, suggested cron snippet.

**Done when:** a fresh VPS can be brought up with `make build`, the Caddy snippet, and the systemd unit — nothing else.

---

## Open questions / future ideas

- Email (SMTP) channel — deferred; webhook covers the MVP.
- Per-URL polling intervals — currently global; revisit if the use case appears.
- Multi-user / role-based auth — explicitly out of scope per CLAUDE.md.
- Status page (public, unauthenticated) — possible follow-up.
