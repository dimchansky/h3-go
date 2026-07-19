# Documentation index

Everything in this directory, what it is for, and how the generated files
stay current. The project front door is the [root README](../README.md);
contributor ground rules are in [CONTRIBUTING.md](../CONTRIBUTING.md), with
a one-page maintainer/agent quick reference in [AGENTS.md](../AGENTS.md).

## Repository map

The root directory is the `h3` library: one flat Go package deliberately
holding the public API layer, the mechanically ported C layer, and the
opt-in parity harness. Why a package split is impossible without breaking
the zero-copy/method architecture — and how to recognize which layer any
root file belongs to — is documented in
[repository-layout-review.md](repository-layout-review.md) (DR-008); the
generated per-file map is
[file-layer-inventory.csv](file-layer-inventory.csv)
(`make layout-inventory`; gated by `make check-layout`).

| Path | What it is |
|---|---|
| [`cmd/h3`](../cmd/h3) | The `h3` executable — a minimal `main` that delegates to `internal/cli` |
| [`internal/cli`](../internal/cli) | CLI implementation: upstream-compatible parser, command registry, output encoders, exit-code mapping; consumes only the public `h3` API |
| [`interop/uberdiff`](../interop/uberdiff) | Separate Go module that differentially tests this library against the official uber/h3-go cgo binding (nightly CI) |
| [`interop/uberbench`](../interop/uberbench) | Separate Go module benchmarking this library against uber/h3-go — equivalence-gated benchmarks plus process-level memory probes; results in [benchmarks](benchmarks/README.md) |
| [`testref`](../testref) | Scaffolding that downloads pristine upstream H3 C sources for the parity suite and gates — never vendored |
| [`tools`](../tools) | Maintenance commands: API/test/CLI inventories, upstream symbol diff, docs link check ([tools/README.md](../tools/README.md)) |
| [`docs`](.) | Design records, compatibility contracts, and generated inventories (this index) |
| [`.github/workflows`](../.github/workflows) | Tiered CI: fast checks + C parity on every code change, heavy suites nightly and on tags ([ci-policy.md](ci-policy.md)) |

## API discovery

| Document | Contents |
|---|---|
| [api-map.md](api-map.md) | Generated C→Go API map: every public H3 C v4.4.0 function with its idiomatic Go equivalent and the additive `Append*`/`*Seq` forms — the quick-discovery projection of the comparison matrix. |

## Design and policy

| Document | Contents |
|---|---|
| [public-api-architecture.md](public-api-architecture.md) | The authoritative design record: API architecture, decision records (DR-001…), measurements, and the living upstream-synchronization workflow (§10). Authored against H3 C v4.3.0 and kept as a snapshot; current facts live in the README/CHANGELOG. |
| [DEVIATIONS.md](DEVIATIONS.md) | Every *intentional* behavioral difference from H3 C. Anything not listed must match C exactly (enforced by the parity suite). Consulted on every upstream sync. |
| [FUTURE_WORK.md](FUTURE_WORK.md) | Deliberately deferred features with full context (GeoJSON, workspaces, …), profiling-gated ideas, rejected designs, and the pre-v1.0.0 checklist. |
| [repository-layout-review.md](repository-layout-review.md) | Why the flat single-package layout stays (DR-008, reaffirming DR-001 with probe evidence), every package-split alternative evaluated and rejected, and the phased discoverability plan behind the file-layer inventory. |
| [versioning.md](versioning.md) | The versioning and release policy: module SemVer and the H3 Core compatibility target as independent version axes, tag rules, release metadata, and the release-note outline. |
| [releasing.md](releasing.md) | The operational release runbook: gates, rc/final-tag procedure, reproducible-artifact verification, immutable Release publication, proxy/pkg.go.dev checks, rulesets, and failure handling. |
| [ci-policy.md](ci-policy.md) | Which CI tier runs when and why the expensive suites are not on every push. |
| [lint-policy.md](lint-policy.md) | Why style-tier lint checks are excluded for mechanically ported files (and only those), the `//nolint` inventory, and when to revisit each exclusion. |

## Comparison with the official cgo binding

| Document | Contents |
|---|---|
| [comparison-uber-h3-go.md](comparison-uber-h3-go.md) | Evidence-based comparison with uber/h3-go: pinned versions, function-by-function coverage matrix (generated), behavioral differences, trade-offs in both directions, and the maintainer checklist for keeping it current. |
| [migration-from-uber-h3-go.md](migration-from-uber-h3-go.md) | Practical migration and upgrade guide: additional capabilities, type/call-site mappings, unit and error-handling changes, and a verified before/after example. |
| [public-api-ergonomics-review.md](public-api-ergonomics-review.md) | Pre-v1 review of every migration/API-shape difference: implement, keep, defer, or reject decisions with compatibility and allocation impact. |
| [benchmarks/README.md](benchmarks/README.md) | Benchmark methodology, memory-accounting caveats, and the committed per-environment result artifacts (raw output, benchstat summaries, process-level memory matrix, environment metadata). |

## CLI

| Document | Contents |
|---|---|
| [cli-compatibility.md](cli-compatibility.md) | The `h3` command-line compatibility contract: command set, formats, exit-code policy, tolerances, and how the inventories below lock it. |
| [../cmd/h3/README.md](../cmd/h3/README.md) | User-facing CLI README: install, usage, version metadata. |

## Upstream test equivalence and syncs

| Document | Contents |
|---|---|
| [ported-c-tests.md](ported-c-tests.md) | Human-readable entry point for upstream test equivalence: audited scope, dispositions, fixture/fuzz suites, and the per-release update procedure. |
| [sync/4.3.0-to-4.4.0.md](sync/4.3.0-to-4.4.0.md) | Reviewed record of the 4.3.0 → 4.4.0 upstream sync — the template future syncs follow. |
| [sync/4.4.0-to-4.5.0.md](sync/4.4.0-to-4.5.0.md) | Discovery record for the 4.4.0 → 4.5.0 upstream migration: full symbol/test/CLI dispositions and the proposed implementation-issue decomposition. Discovery only — the parity target remains 4.4.0 until the implementation issues land. |
| [sync/h3v450-exclusion-inventory.md](sync/h3v450-exclusion-inventory.md) | Temporary migration artifact (#27): every file/test excluded from the `H3VER=4.5.0` parity configuration by the `!h3v450` build tag, mapped to its owning implementation issue. Must be empty (and deleted) at the #36 cutover. |

## Generated inventories (do not edit by hand)

| File | Contents | Regenerate / verify |
|---|---|---|
| [api-surface.txt](api-surface.txt) | Golden lock of the exported Go API surface. | `UPDATE_API_SURFACE=1 go test -run TestAPISurface .` / `make test` |
| [c-api-inventory.csv](c-api-inventory.csv) | Full C-function ↔ Go-declaration mapping. | `make api-inventory` / `make check-api` |
| [file-layer-inventory.csv](file-layer-inventory.csv) | Per-file architectural-layer map of the flat root package (public API vs ported layer vs parity harness). | `make layout-inventory` / `make check-layout` |
| [upstream-test-inventory.csv](upstream-test-inventory.csv) | Case-level registry of every upstream test-ecosystem entry with reviewed dispositions. | reviewed by hand / `make check-test-inventory` |
| [cli-contract.csv](cli-contract.csv) | Semantic command/flag/format contract of the CLI. | reviewed by hand / `make check-cli-inventory` |
| [cli-test-inventory.csv](cli-test-inventory.csv) | All 172 upstream CLI scenarios with expected outputs and source hashes. | `go run ./tools/cliinventory -emit-cases` / `make check-cli-inventory` |
| [cli-fixture-inventory.csv](cli-fixture-inventory.csv) | Upstream CLI input fixtures with hashes. | `go run ./tools/cliinventory -emit-fixtures` / `make check-cli-inventory` |
| [cli-source-inventory.csv](cli-source-inventory.csv) | Upstream sources that define the CLI contract, with hashes. | `go run ./tools/cliinventory -emit-sources` / `make check-cli-inventory` |
| [comparison-uber-h3-go.csv](comparison-uber-h3-go.csv) | Curated per-C-function comparison matrix vs uber/h3-go (source of the generated tables in comparison-uber-h3-go.md and api-map.md). | edit by hand, then `make gen-ubercompare` / `make check-ubercompare` |
| [api-map.md](api-map.md) | Simplified C→Go API map rendered from the comparison matrix (generated section only; the framing text is hand-written). | `make gen-ubercompare` / `make check-ubercompare` |
| [benchmarks/](benchmarks/README.md) | Committed benchmark artifacts per environment (raw, benchstat, memory, metadata), plus a generated all-scenario comparison and README scorecard. | `make bench-uber`; `make gen-benchdocs` / `make check-benchdocs` |

The tools behind these files are documented in
[tools/README.md](../tools/README.md).
