# htmx 4 — Migration cheat-sheet for this project

Source of truth: <https://four.htmx.org/docs/get-started/migration>. This page only captures the bits that affect the **Principled Uptime Tracker** codebase; consult the upstream guide for anything not listed here.

## ⚠️ Beta caveat

The vendored file (`internal/web/static/js/htmx.min.js`) is **htmx `4.0.0-beta2`**. The 4.x API is still stabilising — pin to the file we shipped, not `@latest`. When upgrading:

1. Replace the file in-place under `internal/web/static/js/htmx.min.js`.
2. Bump the version cell in `README.md`'s vendored-asset table.
3. Sweep templates against the rename table below.

## Behavioural defaults that will surprise you

| Default                  | htmx 1.x / 2.x                  | htmx 4.0                             |
| :----------------------- | :------------------------------ | :----------------------------------- |
| Swap on 4xx / 5xx        | **Not swapped**                 | **Swapped** (except 204, 304)        |
| Request timeout          | None (`0`)                      | `60000` ms                           |
| Attribute inheritance    | Implicit, walks ancestors       | **Explicit** — opt in with `:inherited` suffix |

Implications for our app:
- The status-fragment swap "just works" for failed pings — we don't need a custom error-handling path on the client.
- Slow uptime checks still finish well under 60s, so the new default is fine; if a particular monitor needs longer, set it via `<meta name="htmx-config" content='{"defaultTimeout": 90000}'>` or a per-element `hx-config`.
- We don't rely on inheritance anywhere, so nothing to migrate there.

## Renames we will use

### Attributes

| Removed          | Replacement                                       |
| :--------------- | :------------------------------------------------ |
| `hx-vars`        | `hx-vals` (use `js:` prefix when value is JS)     |
| `hx-prompt`      | `hx-confirm` (use `js:` prefix when value is JS)  |
| `hx-ext`         | Include the extension `<script>` tag directly     |
| `hx-disable`     | `hx-ignore` (skip htmx processing)                |
| `hx-disabled-elt`| `hx-disable` (disables form elements during a request — note the swap in semantics) |
| `hx-disinherit`, `hx-inherit` | Removed; inheritance is opt-in via `:inherited` |

### Config

| Old                | New              |
| :----------------- | :--------------- |
| `defaultSwapStyle` | `defaultSwap`    |
| `historyEnabled`   | `history`        |
| `timeout`          | `defaultTimeout` |

Plus 15+ less-used config keys removed (`addedClass`, `allowEval`, `disableSelector`, `historyCacheSize`, `scrollBehavior`, …).

### Events

The convention is now `htmx:phase:action[:sub-action]` with **colons**, not camelCase.

| Old                  | New                  |
| :------------------- | :------------------- |
| `htmx:beforeRequest` | `htmx:before:request`|
| `htmx:afterSwap`     | `htmx:after:swap`    |
| `htmx:beforeOnLoad`  | `htmx:before:init`   |
| `htmx:validation:*`  | Removed              |
| `htmx:xhr:*`         | Removed              |

All error events are consolidated to a single `htmx:error` — listen there if you need a global error hook.

### Headers

| Old (request)         | New                                    |
| :-------------------- | :------------------------------------- |
| `HX-Trigger`          | `HX-Source`, formatted as `tagName#id` |

| Removed (response)            | Notes                                                |
| :---------------------------- | :--------------------------------------------------- |
| `HX-Trigger-After-Swap`       | Use `htmx:after:swap` listener client-side instead.  |
| `HX-Trigger-After-Settle`     | Use `htmx:after:settle`.                             |
| `HX-Prompt`                   | Replaced by `hx-confirm` machinery.                  |

The server now also receives `HX-Request-Type: full | partial` (auto-set by htmx) and an explicit `Accept: text/html`.

### JavaScript API

| Old                       | New / replacement                |
| :------------------------ | :------------------------------- |
| `htmx.addClass(el, c)`    | `el.classList.add(c)`            |
| `htmx.remove(el)`         | `el.remove()`                    |
| `htmx.off(...)`           | `el.removeEventListener(...)`    |
| `htmx.defineExtension`    | `htmx.registerExtension`         |

`htmx.ajax`, `htmx.config`, `htmx.process`, `htmx.trigger` are unchanged.

## New features worth knowing

- **`hx-action` + `hx-method`** — generic alternative to `hx-get` / `hx-post` / etc.
- **`hx-status:NNN`** — per-status-code override (e.g. `hx-status:422='{"target":"#errors"}'` to retarget validation errors).
- **New swap styles** — `innerMorph`, `outerMorph`, `textContent`, `delete`. The morph styles do diff-based DOM updates and preserve focus/state inside the swap region; useful for live-updating dashboards.
- **`<hx-partial>`** — a `<template hx type="partial">` wrapper that names its own swap target; cleaner than scattered `hx-swap-oob` markup.
- **View Transitions API** — opt in with `config.transitions = true` for animated swaps.

## Compat escape hatch

Loading the `htmx-2-compat` extension restores implicit inheritance, the old camelCase event names, and the previous error-default behaviour. We don't need it — this is a greenfield codebase.

## TL;DR for new contributors

- Use kebab-case-with-colons event names.
- Don't expect `:inherited` semantics for free.
- All HTTP responses (except 204/304) become page content; design fragments accordingly.
- Default timeout is 60s.
- When in doubt, search the migration page before reaching for deprecated patterns.
