# Upstream H3 v4.5.0 test equivalence

This is the human-readable entry point for the upstream-test registry. The
case-level source of truth is
[`upstream-test-inventory.csv`](upstream-test-inventory.csv), verified by
`go run ./tools/testinventory -h3ver 4.5.0 -verify` (or
`make check-test-inventory`). Do not maintain a second filename-only list.

## Audited scope

The inventory is discovered from the pristine `testref/h3-4.5.0` tree and
contains 824 reviewed entries:

- 66 `testapps` executables containing 498 named or synthetic cases;
- 172 `h3_bin` CLI cases from 64 registration files;
- 24 fuzzer harnesses;
- 12 benchmark sources;
- 9 filter/CLI executable sources and 2 random-input helpers;
- 10 shared app/test support sources;
- 94 checked-in input fixtures;
- 2 authoritative CMake definitions and the optional generated country
  benchmark pipeline.

Current dispositions are 760 `full`, 1 `partial`, 63 `na`, and zero
`missing` or `deferred` (re-countable at any time with
`make check-test-inventory`, which prints the totals). `full` means the
upstream scenario and assertions were compared with the named Go test, not
merely that its implementation has a cgo parity test.

The one partial entry is `testCellsToLinkedMultiPolygon.c::specificLeak`:
the failure/no-crash regression is ported, while its Valgrind leak-detection
facet is not applicable to garbage-collected Go.

The 172 upstream CLI cases are `full`: the repository ships an
upstream-compatible `h3` executable (`cmd/h3`, implemented in
`internal/cli`), and every registered scenario runs against it in-process
(`TestUpstreamCLICompatibility`), with process-level and opt-in C
differential layers on top — see [cli-compatibility.md](cli-compatibility.md)
for that contract. The remaining `na` rows are legacy stream filters, C
assertion and allocator hooks, random fixture generators, and
performance-only benchmark loops; each has an explicit row-level reason in
the CSV.

## Input-driven and exhaustive suites

Three upstream programs have no `TEST(name)` blocks and consume large golden
datasets. Their pure-Go equivalents are:

- `TestUpstreamLatLngToCellFixtures` for `bc*centers.txt` and
  `rand*centers.txt`;
- `TestUpstreamCellToLatLngFixtures` for `res*ic.txt`;
- `TestUpstreamCellToBoundaryFixtures` for `*cells.txt`.

The fixtures remain in the downloaded reference tree. Run all 526,546
cell/coordinate records (including every boundary vertex comparison) with:

```sh
make -C testref h3-source
make test-upstream-fixtures H3VER=4.5.0
```

CI runs the full fixture suites in the Nightly workflow (as a step of the
`parity` job, reusing its downloaded tree) — nightly, on manual dispatch,
and on every `v*` tag as part of the release gate; see
[ci-policy.md](ci-policy.md). Without `H3_UPSTREAM_FIXTURE_ROOT` these tests
skip, which is why `make test` and the per-push CI path do not execute them.

Exhaustive traversal, local-IJ, directed-edge, vertex, hierarchy, and metric
helpers fail on enumeration/setup errors. They must never silently skip a
base cell or child range, because that would make a passing test compatible
with a truncated domain.

## Fuzzing, parity, and coverage

The four `FuzzUpstream*` targets preserve the raw domains of all 24 upstream
libFuzzer/AFL harnesses: cell/index operations, cell sets, polygons, and
internal IJK coordinates. Smoke them with, for example:

```sh
go test -run '^$' -fuzz FuzzUpstreamCellOperations -fuzztime=10s .
```

These layers answer different questions:

- ported Go tests preserve upstream behavioral scenarios and assertions;
- `cgo && c2go` parity tests compare selected Go calls directly with C;
- fuzzing searches raw input domains for crashes and invariant violations;
- line coverage helps locate weak Go branches, but is not proof of upstream
  behavioral equivalence.

## Updating for a new H3 release

Fetch the new tree, run `upstream-diff`, then run the inventory checker
against the new version. It fails on every newly added/removed case, fixture,
fuzzer, benchmark, filter, helper, or build entry until a reviewer records a
disposition. A `full`, `partial`, or `indirect` row must name a real Go test;
`partial` and `na` require an explanation. `missing` is never accepted by
strict verification.
