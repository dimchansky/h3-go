# Agent and maintainer quick reference

Orientation for coding agents and new maintainers working from a fresh
clone. The authoritative, CI-enforced rules live in
[CONTRIBUTING.md](CONTRIBUTING.md) — read it before changing code. This page
is the map, not the law; if it ever disagrees with CONTRIBUTING.md or the
documents it links, those win.

## What this repository is

A pure-Go, dependency-free port of Uber's H3 C v4.4.0 — all 78 public
functions behind a typed Go API — plus the upstream-compatible `h3` CLI
(`cmd/h3`, implemented in `internal/cli`). The [README](README.md) has the
repository map; [docs/README.md](docs/README.md) indexes every document and
generated inventory; [tools/README.md](tools/README.md) covers the
maintenance tools.

## Invariants you must not break

Details and enforcement in [CONTRIBUTING.md](CONTRIBUTING.md):

- **Production code is pure, safe Go** — no cgo, no `unsafe`, no
  dependencies in any normal build (DR-007; gated by `make check-unsafe`
  and depguard). cgo exists only in the opt-in parity harness behind
  `//go:build cgo && c2go`.
- **Ported files stay C-shaped.** `<cfile>_<function>.go` mirrors one C
  function (`__` marks a C static helper); bodies keep C names and control
  flow, with `// Ported from H3 C: <file>::<name>` attributions that the
  sync tooling parses. Never restructure them for style
  ([docs/lint-policy.md](docs/lint-policy.md) explains which lint checks
  are relaxed there — and only there).
- **C `int` → Go `int32`, C `int64_t` → Go `int64`** (overflow parity).
- **Behavior may differ from C only if listed in
  [docs/DEVIATIONS.md](docs/DEVIATIONS.md).** Everything else must match
  exactly (enforced by the cgo parity suite).
- **The exported API surface is locked** (golden file; regenerate
  deliberately). Public wrappers carry `H3 C API:` doc lines; collection
  APIs keep their zero-allocation `Append*` guarantees.
- **The root package stays flat — do not propose an `internal/` split.**
  Which layer every root file belongs to, why the single package is the
  only legal layout, and the naming rules `make check-layout` enforces are
  in [docs/repository-layout-review.md](docs/repository-layout-review.md)
  (DR-008) and the generated
  [docs/file-layer-inventory.csv](docs/file-layer-inventory.csv).
- **`testref/` is download-only** — upstream C sources are fetched
  (`make -C testref h3-source`), never vendored. Generated inventories
  under `docs/*.csv` are maintained by the tools, not edited by hand.

## Before you push

```sh
make fmt lint test check-unsafe check-layout check-docs
```

plus `make test-c2go` (needs `make -C testref h3-source` once) whenever you
touch ported code or the parity harness. Commit style: focused commits with
area prefixes (`api:`, `c2go:`, `cli:`, `tests:`, `tools:`, `docs:`, `ci:`),
no tool-attribution trailers. Upstream syncs follow the step-by-step
workflow in [CONTRIBUTING.md](CONTRIBUTING.md#porting-a-c-function-upstream-syncs);
CI tiers and what runs when are in [docs/ci-policy.md](docs/ci-policy.md).
