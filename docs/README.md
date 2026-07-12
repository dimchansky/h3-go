# Documentation index

Everything in this directory, what it is for, and how the generated files
stay current. The project front door is the [root README](../README.md);
contributor ground rules are in [CONTRIBUTING.md](../CONTRIBUTING.md).

## Design and policy

| Document | Contents |
|---|---|
| [public-api-architecture.md](public-api-architecture.md) | The authoritative design record: API architecture, decision records (DR-001…), measurements, and the living upstream-synchronization workflow (§10). Authored against H3 C v4.3.0 and kept as a snapshot; current facts live in the README/CHANGELOG. |
| [DEVIATIONS.md](DEVIATIONS.md) | Every *intentional* behavioral difference from H3 C. Anything not listed must match C exactly (enforced by the parity suite). Consulted on every upstream sync. |
| [FUTURE_WORK.md](FUTURE_WORK.md) | Deliberately deferred features with full context (GeoJSON, workspaces, …), profiling-gated ideas, rejected designs, and the pre-v1.0.0 checklist. |
| [ci-policy.md](ci-policy.md) | Which CI tier runs when and why the expensive suites are not on every push. |

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

## Generated inventories (do not edit by hand)

| File | Contents | Regenerate / verify |
|---|---|---|
| [api-surface.txt](api-surface.txt) | Golden lock of the exported Go API surface. | `UPDATE_API_SURFACE=1 go test -run TestAPISurface .` / `make test` |
| [c-api-inventory.csv](c-api-inventory.csv) | Full C-function ↔ Go-declaration mapping. | `make api-inventory` / `make check-api` |
| [upstream-test-inventory.csv](upstream-test-inventory.csv) | Case-level registry of every upstream test-ecosystem entry with reviewed dispositions. | reviewed by hand / `make check-test-inventory` |
| [cli-contract.csv](cli-contract.csv) | Semantic command/flag/format contract of the CLI. | reviewed by hand / `make check-cli-inventory` |
| [cli-test-inventory.csv](cli-test-inventory.csv) | All 170 upstream CLI scenarios with expected outputs and source hashes. | `go run ./tools/cliinventory -emit-cases` / `make check-cli-inventory` |
| [cli-fixture-inventory.csv](cli-fixture-inventory.csv) | Upstream CLI input fixtures with hashes. | `go run ./tools/cliinventory -emit-fixtures` / `make check-cli-inventory` |
| [cli-source-inventory.csv](cli-source-inventory.csv) | Upstream sources that define the CLI contract, with hashes. | `go run ./tools/cliinventory -emit-sources` / `make check-cli-inventory` |

The tools behind these files are documented in
[tools/README.md](../tools/README.md).
