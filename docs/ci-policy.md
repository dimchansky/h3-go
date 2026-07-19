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

Exact allocation budgets are enforced by the ordinary optimized test builds
on both the oldest supported Go version and stable. Race instrumentation
changes escape and allocation profiles, so allocation assertions use a
build-tagged helper: non-race builds measure with `testing.AllocsPerRun`,
while race builds execute each measured closure once for behavioral and race
coverage without applying the production allocation threshold.

## Tiers

| Tier | Jobs | Runs when | Typical duration |
|---|---|---|---|
| Classifier | `changes` | every push/PR | ~10 s |
| Docs | `docs` (Markdown link/anchor gate, `make check-docs`; uber/h3-go comparison-matrix drift gate, `make check-ubercompare`) | every push/PR, including docs-only changes | ~30 s |
| Fast (required) | `fast` (fmt, no-unsafe gate, file-layer gate `make check-layout` (DR-008), lint, smrcptr, pure-Go library + CLI tests and binary build on 1.24 + stable), `api-gates` (API, test, and CLI inventory gates + inventory-freshness diffs) | every push/PR **with code changes** | ~2–3 min to signal |
| Core correctness | `parity` (213-file cgo suite vs original C) | every push/PR with code changes, in parallel with `fast` | ~5–6 min |
| Merge gate | `race` | PRs into the default branch (code changes only) | ~9 min |
| Aggregate gate | `required` (check name **`CI / required`**) — [tools/cirequired](../tools/README.md#cirequired) evaluates every job above against an explicit truth table; the single status check branch protection requires | every push/PR (`if: always()`) | ~30 s after the slowest job |
| Confidence sweep | nightly.yml: `race`, `parity` + gates + the full upstream fixture suites (526,546 golden records), `fuzz-smoke` (6 of the 7 fuzz targets, 60 s each — `FuzzUpstreamPolygonOperations` is seed-corpus-only: issue #3 traced its pathologically slow in-domain inputs to upstream-inherited `maxPolygonToCellsSizeExperimental` behavior, reproduced and slower in H3 C itself, so the target stays out until upstream bounds the estimator), library and CLI C differential suites (uberdiff + uberbench equivalence), CLI cross-builds, `vulncheck` (pinned govulncheck, root module), `secret-scan` (digest-pinned gitleaks over the full history) | nightly 03:17 UTC, `workflow_dispatch`, and every `v*` tag (release gate) | ~15 min wall |
| Release artifacts | release-builds.yml: `build` (single entry point `make release-dist` → tools/releasepack, pinned toolchain, deterministic archives), `verify-reproducible` (independent rebuild; byte-identical SHA256SUMS required), `smoke` (arch-proven runners execute the produced artifacts on linux/amd64, windows/amd64, darwin/arm64) | every `v*` tag (guarded: canonical `vX.Y.Z` / `vX.Y.Z-rc.N` refs only) | ~5 min |
| Benchmarks (informational) | benchmarks.yml: full comparative suite vs the uber/h3-go binding + process-level memory matrix, artifacts uploaded per run | `workflow_dispatch` and monthly schedule — **never a per-push gate** | ~25–35 min |

Notes:

- **Docs-only changes** (`*.md`, `LICENSE`, `NOTICE`, `.gitignore`) skip all
  Go jobs; only the classifier and the docs link check run. Everything else
  — including workflow files, the Makefile, and the generated gate inputs
  `docs/api-surface.txt` / `docs/c-api-inventory.csv` /
  `docs/file-layer-inventory.csv` — counts as code. The
  classifier **fails open**: if it cannot determine the base commit (force
  push, first push), it runs everything.
- **Concurrency cancellation**: a newer commit on the same branch/PR cancels
  in-flight runs of the older one.
- **Direct pushes to master** don't run race; nightly covers them within a
  day, and tags gate releases. Run `go test -race ./...` locally before
  tagging (it is also part of the tag-triggered battery).

## What is required where

- **Every code change**: `fast` + `api-gates` + `parity` green.
- **Before merging a PR**: the above plus `race` — enforced through the
  single required status check **`CI / required`**: an aggregate job that
  runs on every outcome and delegates to
  [tools/cirequired](../tools/README.md#cirequired), which demands the exact
  truth-table row for the event and classifier verdict (docs-only changes
  must *skip* the Go jobs, code changes must *run* them; `race` is expected
  only on PRs). Any failure, cancellation, or unexpected skip — e.g. a
  broken job condition — fails the gate; there is deliberately no generic
  "success or skipped" acceptance. Requiring one aggregate context also
  sidesteps the `race`/`parity` job-name duplication between ci.yml and
  nightly.yml, which would make per-job required contexts ambiguous.
- **Before a release (tag)**: the nightly battery runs on the tag push —
  race, parity + completeness, the upstream fixture suites, fuzz smoke over
  six targets (the seventh is issue #3's documented exception), the
  uber/h3-go differential suite, the pinned
  vulnerability scan, and the full-history secret scan (whose tag-run URL is
  the durable evidence for the released commit). The fixture suites are
  deliberately part of the release gate: they replay every golden
  conversion/boundary record for seconds of compute, so there is no reason
  to release without them. In parallel, release-builds.yml builds the
  archives through the single `make release-dist` entry point, proves
  archive-level reproducibility with an independent rebuild, and
  runtime-smokes the artifacts on architecture-proven runners. The exact
  release procedure is docs/releasing.md; also consult the pre-v1.0.0
  checklist in docs/FUTURE_WORK.md.

Why the fixture suites are nightly rather than per-push: the records are
exercised through the same `LatLngToCell`/`Cell.LatLng`/`Cell.Boundary`
paths that the per-change parity suite already compares against the original
C objects, so a regression that only fixtures would catch is unlikely; and
running them requires the downloaded reference tree, which the fast path
deliberately avoids. They live as a step of the nightly `parity` job to
reuse the tree that job already fetches.
- **Scheduled**: the same battery nightly, catching upstream-source drift
  (parity re-downloads H3 sources) and fuzz regressions.

Why benchmarks are not a CI gate: shared runners have noisy-neighbor
variance well above the deltas that matter, so a pass/fail threshold would
either block unrelated PRs on noise or be too loose to catch real
regressions. The [benchmarks workflow](../.github/workflows/benchmarks.yml)
is manual + monthly, publishes artifacts (raw output, benchstat summaries,
memory matrix, environment metadata), and maintainers promote results into
[docs/benchmarks/](benchmarks/README.md) deliberately — methodology and
noise caveats there.
