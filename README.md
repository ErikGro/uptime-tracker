# Uptime Tracker

A small, self-hosted synthetic monitor for web services. Single Go binary, SQLite, server-rendered HTML. See [`CLAUDE.md`](./CLAUDE.md) for the full spec and [`ROADMAP.md`](./ROADMAP.md) for delivery milestones.

## Configuration

Static configuration (set at deploy time) lives in `.env`:

| Variable        | Purpose                              |
| :-------------- | :----------------------------------- |
| `ADMIN_USER`    | Basic Auth username                  |
| `ADMIN_PASS`    | Basic Auth password                  |
| `LISTEN_ADDR`   | e.g. `:8080`                         |
| `DB_PATH`       | path to the SQLite file              |

Everything else (poll interval, request timeout, failure threshold, webhook URL, retention) is edited at runtime from the dashboard's Settings page.

## Vendored frontend assets

To honour the project's "no CDN, no `node_modules`" rule, frontend libraries are downloaded once and committed under `internal/web/static/`:

| File              | Source                       | Version           |
| :---------------- | :--------------------------- | :---------------- |
| `pico.min.css`    | https://picocss.com          | `2.1.1`           |
| `htmx.min.js`     | https://four.htmx.org        | `4.0.0-beta2`     |
| `alpine.min.js`   | https://alpinejs.dev         | `3.15.11`         |

When upgrading, replace the file in-place and bump the version in this table. For htmx specifically, see [`docs/htmx-4-migration.md`](./docs/htmx-4-migration.md) — htmx 4 is still in beta and has breaking changes vs 1.x/2.x worth knowing.

## Running the POC

```sh
go tool templ generate ./...
go build -o uptime-tracker .
./uptime-tracker
```

Then open <http://localhost:8080/poc> — the page demonstrates the H-A-G-S stack: a Pico-styled layout, an HTMX fragment-swap form, an HTMX polling clock, and an Alpine-driven dialog. The real dashboard at `/` arrives in later milestones; for now `/` redirects to the POC.
