# CI policy

Why expensive jobs do not run on every commit, and what runs when.

## Race-detector value assessment

The production library contains **no goroutines, no `sync`/`atomic` usage,
and no mutable shared state** — its package-level variables are read-only
lookup tables and error sentinels initialized at package load, and every
public API is a pure function over caller-provided data. There are no
caches, pools, or reusable workspaces (deliberately deferred; see
docs/FUTURE_WORK.md item 3). Consequently `go test -race ./...` primarily
guards (a) the test suite's own `t.Parallel()` usage and (b) future
regressions if someone introduces shared state. That is worth keeping — but
not worth ~9 minutes on every push. If a concurrency-bearing feature (e.g.
polyfill workspaces) ever lands, revisit this assessment and move race back
into the per-change tier.

## Tiers

| Tier | Jobs | Runs when | Typical duration |
|---|---|---|---|
| Classifier | `changes` | every push/PR | ~10 s |
| Docs | `docs` (Markdown link/anchor gate, `make check-docs`) | every push/PR, including docs-only changes | ~30 s |
| Fast (required) | `fast` (fmt, no-unsafe gate, lint, smrcptr, pure-Go library + CLI tests and binary build on 1.24 + stable), `api-gates` (API, test, and CLI inventory gates) | every push/PR **with code changes** | ~2–3 min to signal |
| Core correctness | `parity` (227-file cgo suite vs original C) | every push/PR with code changes, in parallel with `fast` | ~5–6 min |
| Merge gate | `race` | PRs into the default branch (code changes only) | ~9 min |
| Confidence sweep | nightly.yml: `race`, `parity` + gates, `fuzz-smoke`, library and CLI C differential suites, CLI cross-builds | nightly 03:17 UTC, `workflow_dispatch`, and every `v*` tag (release gate) | ~15 min wall |

Notes:

- **Docs-only changes** (`*.md`, `LICENSE`, `NOTICE`, `.gitignore`) skip all
  Go jobs; only the classifier and the docs link check run. Everything else
  — including workflow files, the Makefile, and the generated gate inputs
  `docs/api-surface.txt` / `docs/c-api-inventory.csv` — counts as code. The
  classifier **fails open**: if it cannot determine the base commit (force
  push, first push), it runs everything.
- **Concurrency cancellation**: a newer commit on the same branch/PR cancels
  in-flight runs of the older one.
- **Direct pushes to master** don't run race; nightly covers them within a
  day, and tags gate releases. Run `go test -race ./...` locally before
  tagging (it is also part of the tag-triggered battery).

## What is required where

- **Every code change**: `fast` + `api-gates` + `parity` green.
- **Before merging a PR**: the above plus `race`.
- **Before a release (tag)**: the nightly battery runs on the tag push —
  race, parity + completeness, fuzz smoke, and the uber/h3-go differential
  suite. Also consult the pre-v1.0.0 checklist in docs/FUTURE_WORK.md.
- **Scheduled**: the same battery nightly, catching upstream-source drift
  (parity re-downloads H3 sources) and fuzz regressions.
