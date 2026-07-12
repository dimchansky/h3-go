# Lint policy

Why some warnings are deliberately excluded, exactly where the exclusions
apply, and when to revisit them. Configuration lives in
[.golangci.yml](../.golangci.yml); this page is the rationale.

## The two code tiers

1. **Mechanically ported C-shaped code** — root files named
   `<cfile>_<function>.go` / `<cfile>__<staticHelper>.go` and the ported
   upstream test suites (`test<Upstream>…_test.go`, `<cfile>_test…_test.go`).
   These mirror upstream H3 C statement for statement and keep C names
   (CONTRIBUTING.md ground rule 2). Their value is *diffability*: during an
   upstream sync, `tools/upstreamdiff` maps changed C symbols to these files
   and a reviewer compares bodies line by line. Rewriting `x = x + 1` to
   `x++`, converting if/else chains to switch, dropping else branches, or
   renaming `fillIndex_assertions` would make every future sync harder for
   zero behavioral gain.
2. **Idiomatic Go** — the public API files (`cell.go`, `traversal.go`, …),
   `internal/cli`, `cmd/h3`, `tools/`, `interop/uberdiff`, and newly written
   tests. No style exemptions: everything below the "ported tier" carve-out
   applies with full strictness.

Correctness-tier findings (staticcheck, govet, `unused`, `ineffassign`,
depguard, and the non-style gocritic/revive checks) apply to **both** tiers.
Never suppress a correctness, security, or `unsafe` finding to keep ported
code diffable — if C-faithful code genuinely trips one, that is either a
real bug (fix it and document any C divergence in
[DEVIATIONS.md](DEVIATIONS.md)) or a case for a narrowly scoped, explained
`//nolint`.

## Deliberate exclusions (all path-scoped to the ported tier)

| Check | Scope | Why | Revisit when |
|---|---|---|---|
| gocritic `assignOp` | ported files only | C writes `*rotations = *rotations + 1`; the Go port keeps that shape | a file stops being a mechanical port |
| gocritic `ifElseChain` | ported files only | C if/else ladders stay ladders; switch rewrites break line mapping | same |
| revive `var-naming` | ported files only | C identifiers (e.g. `fillIndex_assertions`) are kept verbatim | same |
| revive `var-declaration` | ported files only | C-style explicit zero initializers are kept | same |
| revive `indent-error-flow`, `superfluous-else` | ported files only | C else-branch structure is kept | same |
| revive `increment-decrement` | ported files only | `i += 1` stays when C wrote `i += 1` | same |
| revive `unreachable-code`, `empty-block` | ported files only | defensive C tails/blocks are ported as-is | same |

The path regex in `.golangci.yml` enumerates the C translation-unit
prefixes. A new idiomatic root file will not match it unless named like a
ported file — don't do that. When porting a *new* C file during an upstream
sync, add its prefix to both rules.

Everything else in the ported tier still lints: these files pass
staticcheck, govet, `unused` (see the `//nolint` inventory below), godot,
misspell, and the remaining ~30 gocritic and 15 revive checks.

## Global linter configuration choices

- **depguard bans `unsafe`** everywhere as a lint-time echo of DR-007. The
  authoritative gate is `make check-unsafe` (build-selection analysis across
  all build modes); depguard just fails faster in editors.
- **nolintlint requires every `//nolint` to be specific and explained** —
  a bare `//nolint` fails CI.
- **smrcptr** (via `make lint`) enforces consistent receiver types
  repo-wide; no exclusions.
- **gofmt -s** is a CI gate (`make fmt`); ported code is gofmt-clean — C
  shape is preserved at the statement level, not the whitespace level.

## Removed (stale) exclusions — history

Two global exclusions from the original lint bootstrap were removed in 2026-07
after verifying nothing fires without them:

- `staticcheck SA1019` (use of deprecated APIs) — nothing in the repository
  uses a deprecated API today. If an upstream sync deprecates an H3 function
  that the port must keep calling, re-add this **scoped to the specific
  files**, never globally.
- `govet fieldalignment` in tests — dead rule: the enabled govet analyzer
  set does not include fieldalignment at all.

## `//nolint` directive inventory

Current directives, all in ported files (nolintlint enforces the
reason text):

- `//nolint:unused` on 18 ported utility/debug helpers (`utility_*.go`,
  `vec3d__square.go`, `vertexGraph__initVertexNode.go`,
  `h3index__mode_highbit.go`, `utility_constants.go`): ported for parity
  completeness and exercised only by the `cgo && c2go` parity tests, which
  normal lint builds do not select.
- `//nolint:unparam` on `faceijk__adjustPentVertOverage.go` and
  `coordijk__setIJK.go`: the "redundant" return value / parameter mirrors
  the C signature.

New `//nolint` directives must name the specific linter and carry a reason;
prefer path/text-scoped rules in `.golangci.yml` when a whole tier is
affected.

## Editor-only "modernize" hints

`gopls` suggests modernizations (range-over-`SplitSeq`,
`slices.Contains`, dropping loop-variable copies, `unusedfunc`, …) that are
**not** CI-enforced. Apply them freely in idiomatic code; do **not** apply
them to ported files when they would change statement shape — the hints are
expected background noise there. The `tc := tc` loop-var copies in older
ported tests are harmless and not worth churning.
